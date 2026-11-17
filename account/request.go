package account

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core"
	coreerrors "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core/errors"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/security"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/transport"
)

// RequestOptions carries the optional flags for the encrypted request.
type RequestOptions struct {
	APIVersion string
	Auth       string
	Input      map[string]any
	Method     string
	Encrypt    bool
	TmpSession bool
}

// Request is the central encrypted request dispatcher. It mirrors
// the original Builder: it encrypts the input + sends it through the
// transport layer, decrypts the response, validates `status_det ==
// OK` and returns a typed *events.Event.
func (a *Account) Request(ctx context.Context, opts RequestOptions) (*events.Event, error) {
	// Lazy-init the DC config, auth secret and symmetric key under a
	// lock so concurrent Request calls cannot race on those fields or
	// double-discover DCs.
	a.initMu.Lock()
	if a.dcConfig == nil || a.dcConfig.APIURL == "" {
		cfg, err := transport.DiscoverDCs(ctx, a.transport)
		if err != nil {
			a.initMu.Unlock()
			return nil, err
		}
		a.dcConfig = cfg
	}
	if a.auth == "" {
		secret, err := security.RandomSecret(core.AuthLength)
		if err != nil {
			a.initMu.Unlock()
			return nil, err
		}
		a.auth = secret
	}
	if a.symmetricKey == "" {
		key, err := security.DeriveKey(a.auth)
		if err != nil {
			a.initMu.Unlock()
			return nil, err
		}
		a.symmetricKey = key
	}
	a.initMu.Unlock()

	apiVersion := opts.APIVersion
	if apiVersion == "" {
		apiVersion = core.APIVersion
	}
	auth := opts.Auth
	if auth == "" {
		auth = a.auth
	}
	input := opts.Input
	if input == nil {
		input = map[string]any{}
	}
	method := opts.Method
	if method == "" {
		method = "getUserInfo"
	}
	encrypt := opts.Encrypt
	tmpSession := opts.TmpSession

	// Build the top-level payload.
	payload := map[string]any{"api_version": apiVersion}
	if tmpSession {
		payload["tmp_session"] = auth
	} else {
		// Honour a caller-supplied Auth (RequestOptions.Auth) rather than
		// always falling back to the account's decodedAuth. When opts.Auth
		// is empty, `auth` is a.auth and DecodeAuth(a.auth) == a.decodedAuth,
		// so behaviour is unchanged for the common path.
		payload["auth"] = security.DecodeAuth(auth)
	}

	if apiVersion == "6" {
		encInput := map[string]any{
			"client": a.identity,
			"method": method,
			"input":  input,
		}
		// Verbose debug logging — only when the account was
		// explicitly constructed with debug enabled. The full plaintext
		// is NEVER logged; for sendMessage we emit a non-sensitive
		// summary (chat prefix + text length, not the text itself).
		if a.debug {
			if method == "sendMessage" {
				a.infof("sendMessage input: %s", summarizeSendInput(input))
			} else {
				a.debugf("%s input keys: %v", method, mapKeys(input))
			}
		}
		if encrypt {
			enc, err := security.EncryptJSON(encInput, []byte(a.symmetricKey))
			if err != nil {
				return nil, err
			}
			payload["data_enc"] = enc
			if !tmpSession && a.rsaKey != nil {
				sig, err := security.SignRSA(a.rsaKey, enc)
				if err != nil {
					return nil, err
				}
				payload["sign"] = sig
			}
		} else {
			payload["data_enc"] = encInput
		}
	} else {
		// API versions 0 and 4 (rarely used; parity for completeness).
		switch apiVersion {
		case "0":
			payload["auth"] = auth
			payload["client"] = a.identity
			payload["data"] = input
			payload["method"] = method
		case "4":
			payload["client"] = a.identity
			payload["method"] = method
		}
	}

	body, status, err := a.transport.Post(ctx, a.dcConfig.APIURL, payload)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, coreerrors.New(coreerrors.KindAPI, fmt.Sprintf("http %d: %s", status, string(body)))
	}
	var resp map[string]any
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, coreerrors.Wrap(coreerrors.KindAPI, "decode response", err)
	}
	if enc, ok := resp["data_enc"].(string); ok && enc != "" {
		decrypted, err := security.Decrypt(enc, []byte(a.symmetricKey))
		if err != nil {
			return nil, coreerrors.Wrap(coreerrors.KindAuth, "decrypt", err)
		}
		if err := json.Unmarshal(decrypted, &resp); err != nil {
			return nil, coreerrors.Wrap(coreerrors.KindAPI, "decode decrypted", err)
		}
	}
	st, _ := resp["status"].(string)
	stDet, _ := resp["status_det"].(string)
	if st == "OK" && stDet == "OK" {
		data, _ := resp["data"].(map[string]any)
		if data == nil {
			data = map[string]any{}
		}
		return events.New(data), nil
	}
	// Verbose debug logging of the decrypted server response — only
	// when debug is explicitly enabled. Sensitive fields (auth, sign,
	// data_enc, file_id, access_hash_*, text, ...) are masked.
	if a.debug {
		a.infof("%s server response: %s", method, redactForLog(resp))
	}
	return nil, coreerrors.FromStatus(stDet, fmt.Sprintf("%s failed", method))
}

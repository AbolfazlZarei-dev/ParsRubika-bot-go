# `security` — Rsa

_Source: [`rsa.go`](../../ParsRubika/security/rsa.go)_

## Symbols

### `SignRSA`

```go
func SignRSA(priv *rsa.PrivateKey, data string) (string, error)
```

SignRSA returns the base64-encoded PKCS#1 v1.5 signature of `data`
under the supplied RSA private key. SHA-256 is used as the digest.

### `DecryptOAEP`

```go
func DecryptOAEP(priv *rsa.PrivateKey, data string) (string, error)
```

DecryptOAEP decrypts `data` (base64-encoded) using RSA-OAEP with
SHA-1 (matching PyCryptodome's PKCS1_OAEP default).

### `GenerateKeyPair`

```go
func GenerateKeyPair(bits int) (publicKey string, privateKey string, err error)
```

GenerateKeyPair generates a fresh RSA keypair with the supplied bit
size. It returns:

- publicKey : PEM-encoded RSA public key in SubjectPublicKeyInfo
(PKIX) format — exactly what PyCryptodome's
`key.publickey().export_key()` produces by default
(`-----BEGIN PUBLIC KEY-----`). The Rubika server's RSA encryptor
expects this format; using PKCS#1 (`-----BEGIN RSA PUBLIC KEY-----`)
here causes the server to encrypt the returned auth with a key
derived from the wrong ASN.1 structure, which surfaces as
`security: invalid padding value` during DecryptOAEP.
The PEM is then passed through DecodeAuth.
- privateKey : PEM-encoded PKCS#1 RSA private key, as a string.
- err : any error from key generation / encoding.

### `ParsePrivateKey`

```go
func ParsePrivateKey(s string) (*rsa.PrivateKey, error)
```

ParsePrivateKey parses a PEM-encoded RSA private key (PKCS#1 or
PKCS#8). The input may include the `[REDACTED:...]` marker added by
some clients — it is stripped automatically.

### `ParsePublicKey`

```go
func ParsePublicKey(s string) (*rsa.PublicKey, error)
```

ParsePublicKey parses a PEM-encoded RSA public key (PKCS#1 or PKCS#8).

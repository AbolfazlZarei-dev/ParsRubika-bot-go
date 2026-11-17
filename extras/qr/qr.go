// Package qr provides a simple QR-code generator. It produces a
// PNG image of a QR code from any text/URL, using a pure-Go
// implementation (no external dependencies).
//
// The QR codes are generated at a fixed error-correction level (M)
// and a configurable module size. The output can be saved to disk
// or embedded directly in a message.
package qr

import (
	"bytes"
	"encoding/base64"
	"image/color"
	"os"

	qrcode "github.com/skip2/go-qrcode"
)

// Level is the error-correction level.
type Level int

const (
	LevelL Level = iota // 7% recovery
	LevelM              // 15% recovery
	LevelQ              // 25% recovery
	LevelH              // 30% recovery
)

// Encode generates a QR code PNG image and returns it as a byte
// slice.
func Encode(text string, size int, level Level) ([]byte, error) {
	qr, err := qrcode.New(text, toQrLevel(level))
	if err != nil {
		return nil, err
	}
	qr.ForegroundColor = color.Black
	qr.BackgroundColor = color.White
	pngData, err := qr.PNG(size)
	if err != nil {
		return nil, err
	}
	return pngData, nil
}

// EncodeBase64 generates a QR code and returns it as a base64
// string (useful for embedding in HTML/markdown).
func EncodeBase64(text string, size int, level Level) (string, error) {
	pngData, err := Encode(text, size, level)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(pngData), nil
}

// SaveFile generates a QR code and writes it to the given file.
func SaveFile(text, path string, size int, level Level) error {
	pngData, err := Encode(text, size, level)
	if err != nil {
		return err
	}
	return os.WriteFile(path, pngData, 0o644)
}

// EncodeBuffer generates a QR code and writes it to the supplied
// buffer.
func EncodeBuffer(text string, size int, level Level, buf *bytes.Buffer) error {
	pngData, err := Encode(text, size, level)
	if err != nil {
		return err
	}
	_, err = buf.Write(pngData)
	return err
}

// toQrLevel converts our Level type to the library's level.
func toQrLevel(l Level) qrcode.RecoveryLevel {
	switch l {
	case LevelL:
		return qrcode.Low
	case LevelM:
		return qrcode.Medium
	case LevelQ:
		return qrcode.High
	case LevelH:
		return qrcode.Highest
	}
	return qrcode.Medium
}

// Package models defines the data models used by the bot API client.
// Each model is a plain struct that maps directly to the JSON shape
// returned by the server, so callers can use encoding/json without any
// custom marshaler.
package models

// File is the file-attachment metadata returned by the bot API.
type File struct {
	FileID   string `json:"file_id,omitempty"`
	FileName string `json:"file_name,omitempty"`
	Size     string `json:"size,omitempty"`
}

// Chat represents a chat object returned by `getChat`.
type Chat struct {
	ChatID    string `json:"chat_id"`
	ChatType  string `json:"chat_type"`
	UserID    string `json:"user_id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Title     string `json:"title,omitempty"`
	Username  string `json:"username,omitempty"`
}

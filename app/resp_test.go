package main

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestDecodeBulkString(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		want        BulkString
		wantErr     error
		description string
	}{
		{
			name:        "Valid bulk string",
			input:       []byte("5\r\nhello\r\n"),
			want:        []byte("hello"),
			wantErr:     nil,
			description: "Basic case with a simple string",
		},
		{
			name:        "Empty string",
			input:       []byte("0\r\n\r\n"),
			want:        []byte(""),
			wantErr:     nil,
			description: "Edge case with empty string",
		},
		{
			name:        "Missing CRLF after length",
			input:       []byte("5hello\r\n"),
			want:        nil,
			wantErr:     ErrInvalidEncoding,
			description: "Error case: malformed length delimiter",
		},
		{
			name:        "Invalid length format",
			input:       []byte("abc\r\nhello\r\n"),
			want:        nil,
			wantErr:     ErrInvalidEncoding,
			description: "Error case: non-numeric length",
		},
		{
			name:        "Incomplete input",
			input:       []byte("5\r\nhel"),
			want:        nil,
			wantErr:     io.EOF,
			description: "Error case: truncated input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader(tt.input)
			decoder := &RESPDecoder{Reader: reader}

			got, err := decoder.Decode()

			// Check error
			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("DecodeBulkString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// For specific error types, check if they match
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("DecodeBulkString() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			gotBytes, ok := got.([]byte)
			if !ok {
				t.Errorf("DecodeBulkString() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Check result
			if !bytes.Equal(gotBytes, tt.want) {
				t.Errorf("DecodeBulkString() = %v, want %v", got, tt.want)
			}
		})
	}
}

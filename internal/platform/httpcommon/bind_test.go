package httpcommon

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeJSONRejectsUnknownFields(t *testing.T) {
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"email":"user@example.com","unknown":"x"}`))

	var dst struct {
		Email string `json:"email"`
	}

	err := DecodeJSON(req, &dst)
	if err != ErrInvalidJSON {
		t.Fatalf("expected ErrInvalidJSON, got %v", err)
	}
}

func TestDecodeJSONRejectsMultipleJSONObjects(t *testing.T) {
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"email":"user@example.com"}{"email":"second@example.com"}`))

	var dst struct {
		Email string `json:"email"`
	}

	err := DecodeJSON(req, &dst)
	if err != ErrInvalidJSON {
		t.Fatalf("expected ErrInvalidJSON, got %v", err)
	}
}

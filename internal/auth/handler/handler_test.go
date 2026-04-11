package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ChernykhITMO/Wishlist-API/internal/auth/domain"
	"github.com/ChernykhITMO/Wishlist-API/internal/auth/services"
	"github.com/google/uuid"
)

type userRepositoryStub struct {
	createFn func(ctx context.Context, id uuid.UUID, email, password string) error
}

func (s userRepositoryStub) Create(ctx context.Context, id uuid.UUID, email, password string) error {
	if s.createFn != nil {
		return s.createFn(ctx, id, email, password)
	}
	return nil
}

func (s userRepositoryStub) GetByEmail(ctx context.Context, email string) (uuid.UUID, string, error) {
	return uuid.Nil, "", errors.New("not implemented")
}

type tokenManagerStub struct{}

func (tokenManagerStub) Issue(userID uuid.UUID) (string, error) {
	return "issued-token", nil
}

type passwordManagerStub struct{}

func (passwordManagerStub) Hash(password string) (string, error) {
	return "hashed-password", nil
}

func (passwordManagerStub) Compare(hash, password string) error {
	return nil
}

func TestRegisterReturnsConflictWhenEmailAlreadyExists(t *testing.T) {
	svc := services.New(
		userRepositoryStub{
			createFn: func(ctx context.Context, id uuid.UUID, email, password string) error {
				return domain.ErrEmailAlreadyExists
			},
		},
		tokenManagerStub{},
		passwordManagerStub{},
	)

	h := New(*svc)
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"email":"user@example.com","password":"secret123"}`))
	rec := httptest.NewRecorder()

	h.Register(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, rec.Code)
	}

	var body struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if body.Error.Code != "CONFLICT" {
		t.Fatalf("expected error code CONFLICT, got %q", body.Error.Code)
	}
	if body.Error.Message != "email already exists" {
		t.Fatalf("expected error message %q, got %q", "email already exists", body.Error.Message)
	}
}

func TestRegisterReturnsCreatedWithToken(t *testing.T) {
	svc := services.New(
		userRepositoryStub{},
		tokenManagerStub{},
		passwordManagerStub{},
	)

	h := New(*svc)
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"email":"user@example.com","password":"secret123"}`))
	rec := httptest.NewRecorder()

	h.Register(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var body RegisterResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if body.Token != "issued-token" {
		t.Fatalf("expected token %q, got %q", "issued-token", body.Token)
	}
}

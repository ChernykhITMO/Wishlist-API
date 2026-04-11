package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/ChernykhITMO/Wishlist-API/internal/auth/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"net/mail"
	"strings"
)

type Service struct {
	userRepository UserRepository
	tokenManager   TokenManager
	passwords      PasswordManager
}

func New(
	userRepository UserRepository,
	tokenManager TokenManager,
	passwords PasswordManager) *Service {
	return &Service{
		userRepository: userRepository,
		tokenManager:   tokenManager,
		passwords:      passwords,
	}
}
func (s *Service) Register(ctx context.Context, input RegisterInput) (string, error) {
	id := uuid.New()

	email := strings.TrimSpace(strings.ToLower(input.Email))
	if _, err := mail.ParseAddress(email); err != nil {
		return "", domain.ErrInvalidRegistration
	}

	if strings.TrimSpace(input.Password) == "" {
		return "", domain.ErrInvalidRegistration
	}

	passwordHash, err := s.passwords.Hash(input.Password)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	err = s.userRepository.Create(ctx, id, email, passwordHash)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			return "", domain.ErrEmailAlreadyExists
		}

		return "", fmt.Errorf("create user: %w", err)
	}

	token, err := s.tokenManager.Issue(id)
	if err != nil {
		return "", fmt.Errorf("issue token: %w", err)
	}

	return token, nil
}

func (s *Service) Login(ctx context.Context, input LoginInput) (string, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))
	if _, err := mail.ParseAddress(email); err != nil {
		return "", domain.ErrInvalidCredentials
	}
	id, pass, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", domain.ErrInvalidCredentials
		}
		return "", fmt.Errorf("get user by email: %w", err)
	}

	if err := s.passwords.Compare(pass, input.Password); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	token, err := s.tokenManager.Issue(id)
	if err != nil {
		return "", fmt.Errorf("issue token: %w", err)
	}

	return token, nil
}

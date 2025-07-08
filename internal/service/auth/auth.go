package auth

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	svcErr "github.com/passwordhash/jwt-test-task/internal/service/errors"
)

type JWTService interface {
	NewToken(alg string, sub any, ttl time.Duration, secret string) (string, error)
}

type Service struct {
	log        *slog.Logger
	jwtService JWTService
}

func NewService(
	log *slog.Logger,
	jwtService JWTService,
) *Service {
	return &Service{
		log:        log,
		jwtService: jwtService,
	}
}

func (s *Service) GetPair(ctx context.Context, id, ip, userAgent string) (access, refresh string, err error) {
	const op = "auth.service.GetPair"

	log := s.log.With("op", op, "id", id, "ip", ip, "userAgent", userAgent)

	if _, err := uuid.Parse(id); err != nil {
		log.Warn("invalid id", slog.String("id", id), slog.Any("error", err))

		return "", "", svcErr.ErrInvalidID
	}

	const secret = "some_strong_secret"
	access, err = s.jwtService.NewToken("HS512", id, 15*time.Minute, secret)
	if err != nil {
		log.Error("failed to create access token", slog.Any("error", err))

		return "", "", err
	}

	return access, "bar", nil
}

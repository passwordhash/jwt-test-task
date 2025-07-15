package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/google/uuid"

	svcErr "github.com/passwordhash/jwt-test-task/internal/service/errors"
	"github.com/passwordhash/jwt-test-task/pkg/jwt"
)

const (
	refreshTokenLength = 32
)

const (
	claimTokenID = "token_id"
)

type RefreshTokenSaver interface {
	Save(ctx context.Context, userID, tokenID, tokenHash, userAgent, ip string) (string, error)
}

type RefreshTokenRevoker interface {
	Revoke(ctx context.Context, userID, userAgent string) error
}

type RefreshTokenGenerator interface {
	Generate(length int) (string, error)
	Hash(token string) (string, error)
}

type Service struct {
	log *slog.Logger

	refreshSaver          RefreshTokenSaver
	refreshTokenGenerator RefreshTokenGenerator
	refreshTokenRevoker   RefreshTokenRevoker

	accessTTL time.Duration
	secret    string
}

func New(
	log *slog.Logger,

	refreshSaver RefreshTokenSaver,
	refreshTokenGenerator RefreshTokenGenerator,
	refreshTokenRevoker RefreshTokenRevoker,

	accessTTL time.Duration,
	secret string,
) *Service {
	return &Service{
		log:                   log,
		refreshSaver:          refreshSaver,
		refreshTokenGenerator: refreshTokenGenerator,
		refreshTokenRevoker:   refreshTokenRevoker,
		accessTTL:             accessTTL,
		secret:                secret,
	}
}

func (s *Service) GetPair(
	ctx context.Context,
	userID, remoteAddr, userAgent string,
) (access, refresh string, err error) {
	const op = "tokens.service.GetPair"

	log := s.log.With("op", op, "userID", userID, "remoteAddr", remoteAddr, "userAgent", userAgent)

	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		log.Error("failed to split host and port from IP", slog.Any("error", err))

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	if _, err := uuid.Parse(userID); err != nil {
		log.Warn("invalid uuid format of userID", slog.Any("error", err))

		return "", "", svcErr.ErrInvalidID
	}

	tokenID := uuid.NewString()
	claims := map[string]any{
		"sub":        userID,
		claimTokenID: tokenID,
	}

	access, err = jwt.NewToken("HS512", claims, s.accessTTL, s.secret)
	if err != nil {
		log.Error("failed to create access token", slog.Any("error", err))

		return "", "", err
	}

	refreshHash, err := s.refreshTokenGenerator.Generate(refreshTokenLength)
	if err != nil {
		log.Error("failed to generate refresh token", slog.Any("error", err))

		return "", "", err
	}

	_, err = s.refreshSaver.Save(ctx, userID, tokenID, refreshHash, userAgent, ip)
	// TODO: handle error properly
	if err != nil {
		log.Error("failed to save refresh token", slog.Any("error", err))

		return "", "", err
	}

	log.Info("tokens generated successfully")

	return access, refresh, nil
}

func (s *Service) UserIDByToken(_ context.Context, token string) (string, error) {
	const op = "tokens.service.UserIDByToken"

	log := s.log.With("op", op, "token", token)

	claims, err := jwt.ParseToken(token, s.secret)
	if err != nil {
		log.Error("failed to get user ID from token", slog.Any("error", err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	userID := claims["sub"]

	log.Info("user ID from token", slog.Any("userID", userID))

	return fmt.Sprintf("%v", userID), nil
}

func (s *Service) RevokeRefreshToken(ctx context.Context, userID, userAgent string) error {
	const op = "tokens.service.RevokeRefreshToken"

	log := s.log.With("op", op, "userID", userID, "userAgent", userAgent)

	if _, err := uuid.Parse(userID); err != nil {
		log.Warn("invalid user ID", slog.String("userID", userID), slog.Any("error", err))

		return svcErr.ErrInvalidID
	}

	err := s.refreshTokenRevoker.Revoke(ctx, userID, userAgent)
	if err != nil {
		log.Error("failed to revoke refresh token", slog.Any("error", err))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("refresh token revoked successfully")

	return nil
}

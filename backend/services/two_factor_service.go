package services

import "context"

type UserRepository interface {
	Save2FASecret(ctx context.Context, userId, secret string) error
	Enable2FA(ctx context.Context, userId string) error
	Disable2FA(ctx context.Context, userId string) error
}

type TwoFactorService struct {
	Repo UserRepository
}

func NewTwoFactorService(repo UserRepository) *TwoFactorService {
	return &TwoFactorService{Repo: repo}
}

func (s *TwoFactorService) Generate2FA(userId string) {}

func (s *TwoFactorService) Enable2FA(userId, code string) {}

func (s *TwoFactorService) Disable2FA(userId, password string) {}

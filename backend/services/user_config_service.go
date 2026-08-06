package services

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
)

type UserConfigRepo interface {
	GetUserByIdForAuth(ctx context.Context, id string) (*models.User, error)
	ChangePassword(ctx context.Context, id, newPassword string) error
	DefinePassword(ctx context.Context, id, newPassword string) error
	GetUserByEmailForAuth(ctx context.Context, email string) (*models.User, error)
	UpdateRecoveryToken(ctx context.Context, userId, tokenHash string, expiresAt time.Time) error
	GetUserByRecoveryToken(ctx context.Context, tokenHash string) (*models.User, error)
	GetCredentialsUser(ctx context.Context, id string) (*models.User, error)
}

type AuthConfigRepo interface {
	DeleteOtherFamilies(ctx context.Context, userId, currentFamilyId string) error
	DeleteAllUserRefreshTokens(ctx context.Context, userId string) error
	GetRefreshTokenById(ctx context.Context, id string) (*models.RefreshToken, error)
}

type UserConfigService struct {
	UserRepo         UserConfigRepo
	AuthRepo         AuthConfigRepo
	EmailService     EmailService
	ArgonParams      utils.ArgonParams
	JwtRefreshSecret string
	AppDomain        string
}

func NewUserConfigService(userRepo UserConfigRepo, authRepo AuthConfigRepo, emailService EmailService, argonParams utils.ArgonParams, jwtRSecret, appDomain string) *UserConfigService {
	return &UserConfigService{
		UserRepo:         userRepo,
		ArgonParams:      argonParams,
		EmailService:     emailService,
		AuthRepo:         authRepo,
		JwtRefreshSecret: jwtRSecret,
		AppDomain:        appDomain,
	}
}

func (s *UserConfigService) ChangePassword(ctx context.Context, userIdContext, refreshTokenString string, data dto.ChangePasswordRequest) error {
	if data.NewPassword != data.ConfirmNewPassword {
		return models.ErrPasswordsDontMatch
	}

	user, token, err := s.validateUserSession(ctx, userIdContext, refreshTokenString)
	if err != nil {
		return err
	}

	if user.PasswordHash == nil {
		return models.ErrPasswordNotSet
	}

	ok, err := utils.VerifyPassword(data.OldPassword, *user.PasswordHash)
	if err != nil {
		slog.ErrorContext(ctx, "Argon2 password verification failed during password change",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return models.ErrPasswordVerificationFailed
	}
	if !ok {
		return models.ErrIncorrectPassword
	}

	newPasswordHash, err := utils.HashPassword(data.NewPassword, s.ArgonParams)
	if err != nil {
		slog.ErrorContext(ctx, "failed to generate Argon2 hash for new password",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return models.ErrPasswordChangeFailed
	}

	err = s.UserRepo.ChangePassword(ctx, user.Id, newPasswordHash)
	if err != nil {
		slog.ErrorContext(ctx, "database error changing password",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return models.ErrPasswordChangeFailed
	}

	err = s.AuthRepo.DeleteOtherFamilies(ctx, user.Id, token.FamilyId)
	if err != nil {
		slog.ErrorContext(ctx, "failed to delete other session families after password change",
			slog.String("user_id", user.Id),
			slog.String("current_family_id", token.FamilyId),
			slog.Any("error", err),
		)
		return models.ErrPasswordChangeButNotLogout
	}

	return nil
}

func (s *UserConfigService) DefinePassword(ctx context.Context, userIdContext, refreshTokenString string, data dto.DefinePasswordRequest) error {
	if data.NewPassword != data.ConfirmNewPassword {
		return models.ErrPasswordsDontMatch
	}

	if !utils.IsPasswordValid(data.NewPassword) {
		return models.ErrPasswordIsNotValid
	}

	user, token, err := s.validateUserSession(ctx, userIdContext, refreshTokenString)
	if err != nil {
		return err
	}

	if user.PasswordHash != nil {
		return models.ErrPasswordSet
	}

	newPasswordHash, err := utils.HashPassword(data.NewPassword, s.ArgonParams)
	if err != nil {
		slog.ErrorContext(ctx, "failed to hash new defined password",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return models.ErrPasswordSetFailed
	}

	err = s.UserRepo.DefinePassword(ctx, user.Id, newPasswordHash)
	if err != nil {
		slog.ErrorContext(ctx, "database error defining initial password",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return models.ErrPasswordSetFailed
	}

	err = s.AuthRepo.DeleteOtherFamilies(ctx, user.Id, token.FamilyId)
	if err != nil {
		slog.ErrorContext(ctx, "failed to delete other session families after password definition",
			slog.String("user_id", user.Id),
			slog.String("current_family_id", token.FamilyId),
			slog.Any("error", err),
		)
		return models.ErrPasswordSetButNotLogout
	}

	return nil
}

func (s *UserConfigService) ForgotPassword(ctx context.Context, email string) error {
	maskedEmail := utils.MaskEmail(email)
	user, err := s.UserRepo.GetUserByEmailForAuth(ctx, email)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil
		}
		slog.ErrorContext(ctx, "database error fetching user by email during forgot password",
			slog.String("email", maskedEmail),
			slog.Any("error", err),
		)
		return nil
	}

	if user == nil || !user.Enable {
		return nil
	}

	if user.RecoveryTokenExpiresAt != nil {
		now := time.Now()

		if now.Before(*user.RecoveryTokenExpiresAt) {
			minTimeRemainingForNewSend := 5 * time.Minute
			timeLeft := user.RecoveryTokenExpiresAt.Sub(now)

			if timeLeft > minTimeRemainingForNewSend {
				waitTime := timeLeft - minTimeRemainingForNewSend
				slog.WarnContext(ctx, "forgot password email blocked by cooldown",
					slog.String("email", maskedEmail),
					slog.Duration("wait_time", waitTime.Round(time.Second)),
				)
				return nil
			}
		}
	}

	token, err := utils.GenerateCustomId(32)
	if err != nil {
		slog.ErrorContext(ctx, "failed to generate secure recovery token",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return models.ErrGeneratingToken
	}

	tokenHash := utils.HashSHA512(token)
	expiresAt := time.Now().Add(15 * time.Minute)

	err = s.UserRepo.UpdateRecoveryToken(ctx, user.Id, tokenHash, expiresAt)
	if err != nil {
		slog.ErrorContext(ctx, "failed to save recovery token in database",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return models.ErrGeneratingToken
	}

	asyncCtx := context.WithoutCancel(ctx)
	utils.GoSafe(func() {
		err := s.EmailService.SendRecoveryEmail(user.Email, token)
		if err != nil {
			slog.ErrorContext(asyncCtx, "failed to send recovery email in background",
				slog.String("user_id", user.Id),
				slog.Any("error", err),
			)
		}
	})

	return nil
}

func (s *UserConfigService) ResetPassword(ctx context.Context, data dto.ResetPasswordRequest) error {
	if data.NewPassword != data.ConfirmNewPassword {
		return models.ErrPasswordsDontMatch
	}

	if !utils.IsPasswordValid(data.NewPassword) {
		return models.ErrPasswordIsNotValid
	}

	hashToken := utils.HashSHA512(data.Token)
	user, err := s.UserRepo.GetUserByRecoveryToken(ctx, hashToken)

	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) || errors.Is(err, pgx.ErrNoRows) {
			return models.ErrInvalidOrExpiredToken
		}
		slog.ErrorContext(ctx, "database error looking up recovery token during reset password", slog.Any("error", err))
		return models.ErrInvalidOrExpiredToken
	}

	if user == nil || user.RecoveryTokenExpiresAt == nil || !time.Now().Before(*user.RecoveryTokenExpiresAt) {
		return models.ErrInvalidOrExpiredToken
	}

	if !user.Enable {
		return models.ErrUserNotEnabled
	}

	newPasswordHash, err := utils.HashPassword(data.NewPassword, s.ArgonParams)
	if err != nil {
		slog.ErrorContext(ctx, "failed to hash password for reset password",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return models.ErrPasswordChangeFailed
	}

	err = s.UserRepo.ChangePassword(ctx, user.Id, newPasswordHash)
	if err != nil {
		slog.ErrorContext(ctx, "database error updating password during reset",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return models.ErrPasswordChangeFailed
	}

	err = s.AuthRepo.DeleteAllUserRefreshTokens(ctx, user.Id)
	if err != nil {
		slog.ErrorContext(ctx, "failed to revoke all refresh tokens after password reset",
			slog.String("user_id", user.Id),
			slog.Any("error", err),
		)
		return models.ErrPasswordChangeButNotLogout
	}

	return nil
}

func (s *UserConfigService) GetMyCredentials(ctx context.Context, id string) (*dto.GetMyCredentialsResponse, error) {
	user, err := s.UserRepo.GetCredentialsUser(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, models.ErrUserNotFound
		}
		slog.ErrorContext(ctx, "database query error fetching credentials by user ID",
			slog.String("user_id", id),
			slog.Any("error", err),
		)
		return nil, models.ErrFailQueryUser
	}

	response := &dto.GetMyCredentialsResponse{
		Id:        user.Id,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		RoleId:    user.RoleId,
		LastLogin: user.LastLogin,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return response, nil
}

func (s *UserConfigService) validateUserSession(ctx context.Context, userIdContext, refreshTokenString string) (*models.User, *models.RefreshToken, error) {
	claims, err := utils.ValidateToken(refreshTokenString, s.JwtRefreshSecret, s.AppDomain)
	if err != nil {
		slog.WarnContext(ctx, "validateUserSession: JWT validation failed", slog.Any("error", err))
		return nil, nil, models.ErrInvalidOrExpiredRefresh
	}

	tokenExists, err := s.AuthRepo.GetRefreshTokenById(ctx, claims.ID)
	if err != nil || tokenExists == nil {
		slog.WarnContext(ctx, "validateUserSession: refresh token not found in database",
			slog.String("token_id", claims.ID),
			slog.Any("error", err),
		)
		return nil, nil, models.ErrInvalidOrExpiredRefresh
	}

	if tokenExists.UserId != claims.Subject || userIdContext != claims.Subject {
		slog.WarnContext(ctx, "validateUserSession: metadata mismatch",
			slog.String("context_user_id", userIdContext),
			slog.String("token_subject", claims.Subject),
			slog.String("db_user_id", tokenExists.UserId),
		)
		return nil, nil, models.ErrTokenMetadataMisMatch
	}

	user, err := s.UserRepo.GetUserByIdForAuth(ctx, claims.Subject)
	if err != nil {
		slog.ErrorContext(ctx, "validateUserSession: failed to query user in database",
			slog.String("user_id", claims.Subject),
			slog.Any("error", err),
		)
		return nil, nil, models.ErrUserNotFound
	}

	return user, tokenExists, nil
}

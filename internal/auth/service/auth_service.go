package service

import (
	"context"
	"kinoza-back/internal/auth/domain"
	"kinoza-back/pkg/logger"
	"kinoza-back/pkg/transaction"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var _ domain.UserService = (*UserServiceImpl)(nil)

type UserServiceImpl struct {
	us         domain.UserRepo
	rt         domain.RefreshTokens
	secret     string
	accessTTL  time.Duration
	refreshTTL time.Duration
	tx         *transaction.Transactor
}

func NewUserService(us domain.UserRepo, rt domain.RefreshTokens, secret string, accessTTL, refreshTTL time.Duration, tx *transaction.Transactor) *UserServiceImpl {
	return &UserServiceImpl{us: us, rt: rt, secret: secret, accessTTL: accessTTL, refreshTTL: refreshTTL, tx: tx}
}

// Register - валидирует email и role юзера, хэширует пароль через bcrypt
// (сырой пароль в структуре не попадает наружу), сохраняет нового
// юзера через репозиторий
func (s *UserServiceImpl) Register(ctx context.Context, u *domain.User) (*domain.User, error) {
	if err := u.EmailWrapper(); err != nil {
		return nil, logger.LogErr("failed to validate user", err)
	}
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), 12)
	if err != nil {
		return nil, logger.LogErr("failed to hash password", err)
	}
	u.PasswordHash = string(hashedPw)
	u.Role = "viewer"
	user, err := s.us.CreateUser(ctx, u)
	if err != nil {
		return nil, logger.LogErr("failed to create user", err)
	}
	return user, nil
}

// Login - находит юзера по email, сверяет пароль с хэшем в БД, обе ветки
// отказа (юзер не найден, пароль не совпал) отдают наружу один и тот же
// ErrInvalidCredentials, чтобы не давать возможность перебором узнавать
// зарегистрированные email, при успехе выдаёт пару токенов: access (JWT,
// коротко живущий, подписан) и refresh (случайная строка; в БД сохраняется
// только её хэш, сырое значение уходит клиенту и нигде больше не хранится)
func (s *UserServiceImpl) Login(ctx context.Context, email, password string) (access string, refresh string, err error) {
	user, err := s.us.GetUserByEmail(ctx, email)
	if err != nil {
		slog.Error("failed to get user by email", "error", err)
		return "", "", domain.ErrInvalidCredentials
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		slog.Error("failed to compare password", "error", err)
		return "", "", domain.ErrInvalidCredentials
	}
	access, err = newAccessToken(user.Id, email, user.Role, s.secret, s.accessTTL)
	if err != nil {
		return "", "", logger.LogErr("failed to create access token", err)
	}
	refresh, hash, err := newRefreshToken()
	if err != nil {
		return "", "", logger.LogErr("failed to create refresh token", err)
	}
	token := &domain.RefreshToken{
		UserId:    user.Id,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(s.refreshTTL),
	}
	_, err = s.rt.CreateToken(ctx, token)
	if err != nil {
		return "", "", logger.LogErr("failed to create refresh token", err)
	}
	return access, refresh, nil
}

func (s *UserServiceImpl) GetUserById(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.us.GetUserById(ctx, id)
	if err != nil {
		return nil, logger.LogErr("failed to get user by id", err)
	}
	return user, nil
}

// TODO: фильтрация по role
func (s *UserServiceImpl) UpdateUser(ctx context.Context, u *domain.User, name, email string) error {
	if err := s.us.UpdateUser(ctx, u, name, email); err != nil {
		return logger.LogErr("failed to update user", err)
	}
	return nil
}

// Refresh - хэширует присланный клиентом raw-токен и ищет запись в БД по
// этому хэшу (не по сырому значению - оно там не хранится), отклоняет,
// если токен не найден, уже отозван или истёк, происходитт ротация: найденный токен
// всегда отзывается, даже при успехе - использованный refresh-токен
// одноразовый, подтягивает актуального юзера (email/role могли смениться
// с момента выдачи старого токена) и выдаёт новую пару access/refresh
func (s *UserServiceImpl) Refresh(ctx context.Context, rawToken string) (access, refresh string, err error) {
	hash := hashToken(rawToken)
	token, err := s.rt.GetTokenByHash(ctx, hash)
	if err != nil {
		return "", "", logger.LogErr("failed to get token by hash", err)
	}
	if token.RevokedAt != nil {
		slog.Error("token already revoked", "id", token.Id)
		return "", "", domain.ErrTokenRevoked
	}
	if token.ExpiresAt.Before(time.Now()) {
		slog.Error("token expired", "id", token.Id)
		return "", "", domain.ErrTokenExpired
	}
	user, err := s.us.GetUserById(ctx, token.UserId)
	if err != nil {
		return "", "", logger.LogErr("failed to get user by id", err)
	}
	access, err = newAccessToken(user.Id, user.Email, user.Role, s.secret, s.accessTTL)
	if err != nil {
		return "", "", logger.LogErr("failed to create access token", err)
	}
	newRaw, newHash, err := newRefreshToken()
	if err != nil {
		return "", "", logger.LogErr("failed to create refresh token", err)
	}
	if err = s.tx.Transaction(ctx, func(ctx context.Context) error {
		if err := s.rt.RevokeToken(ctx, token.Id); err != nil {
			return err
		}
		_, err := s.rt.CreateToken(ctx, &domain.RefreshToken{
			UserId:    user.Id,
			TokenHash: newHash,
			ExpiresAt: time.Now().Add(s.refreshTTL),
		})
		return err
	}); err != nil {
		return "", "", logger.LogErr("failed to refresh token", err)
	}
	return access, newRaw, nil
}

func (s *UserServiceImpl) GetTokensList(ctx context.Context, userId string) ([]*domain.RefreshToken, error) {
	tokens, err := s.rt.GetTokensList(ctx, userId)
	if err != nil {
		return nil, logger.LogErr("failed to get tokens list", err)
	}
	return tokens, nil
}

func (s *UserServiceImpl) RevokeToken(ctx context.Context, id string) error {
	if err := s.rt.RevokeToken(ctx, id); err != nil {
		return logger.LogErr("failed to revoke token", err)
	}
	return nil
}

func (s *UserServiceImpl) RevokeAllTokens(ctx context.Context, userId string) error {
	if err := s.rt.RevokeAllTokens(ctx, userId); err != nil {
		return logger.LogErr("failed to revoke all tokens", err)
	}
	return nil
}

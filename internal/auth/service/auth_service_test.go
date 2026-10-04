package service

import (
	"context"
	"errors"
	"kinoza-back/internal/auth/domain"
	"kinoza-back/internal/auth/mocks"
	"kinoza-back/pkg/transaction"
	"log/slog"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func Test_RegisterSuccess(t *testing.T) {
	userMock := new(mocks.UserMock)
	tokenMock := new(mocks.TokenMock)
	model := &domain.User{Email: "test@example.com"}
	accessTTL := 24 * time.Hour
	refreshTTL := 168 * time.Hour
	tx := &transaction.Transactor{}
	svc := NewUserService(userMock, tokenMock, "test_secret", accessTTL, refreshTTL, tx)
	userMock.On("CreateUser", mock.Anything, mock.Anything).Return(model, nil)

	_, err := svc.Register(context.Background(), model)

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	userMock.AssertExpectations(t)
}

func Test_RegisterFail(t *testing.T) {
	userMock := new(mocks.UserMock)
	tokenMock := new(mocks.TokenMock)
	model := &domain.User{Email: "test@example.com"}
	accessTTL := 24 * time.Hour
	refreshTTL := 168 * time.Hour
	tx := &transaction.Transactor{}
	svc := NewUserService(userMock, tokenMock, "test_secret", accessTTL, refreshTTL, tx)
	userMock.On("CreateUser", mock.Anything, mock.Anything).Return(model, errors.New("failed register"))

	_, err := svc.Register(context.Background(), model)

	assert.Equal(t, errors.New("failed register"), err)
	userMock.AssertExpectations(t)
}

func Test_RegisterInvalidEmail(t *testing.T) {
	userMock := new(mocks.UserMock)
	tokenMock := new(mocks.TokenMock)
	model := &domain.User{Email: "testexample.com"}
	accessTTL := 24 * time.Hour
	refreshTTL := 168 * time.Hour
	tx := &transaction.Transactor{}
	svc := NewUserService(userMock, tokenMock, "test_secret", accessTTL, refreshTTL, tx)

	_, err := svc.Register(context.Background(), model)

	assert.Equal(t, domain.ErrInvalidEmail, err)
	assert.Error(t, err)
	userMock.AssertNotCalled(t, "CreateUser", mock.Anything, mock.Anything)
}

func Test_LoginSuccess(t *testing.T) {
	userMock := new(mocks.UserMock)
	tokenMock := new(mocks.TokenMock)
	accessTTL := 24 * time.Hour
	refreshTTL := 168 * time.Hour
	tx := &transaction.Transactor{}
	email := "test@example.com"
	pass := "123qwe"
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
	}
	svc := NewUserService(userMock, tokenMock, "test_secret", accessTTL, refreshTTL, tx)
	userMock.On("GetUserByEmail", mock.Anything, mock.Anything).Return(&domain.User{PasswordHash: string(hash)}, nil)
	tokenMock.On("CreateToken", mock.Anything, mock.Anything).Return(&domain.RefreshToken{}, nil)

	access, refresh, err := svc.Login(context.Background(), email, pass)

	assert.NotEmpty(t, access)
	assert.NotEmpty(t, refresh)
	assert.Equal(t, nil, err)
	userMock.AssertExpectations(t)
	tokenMock.AssertExpectations(t)
}

func Test_LoginFail(t *testing.T) {
	userMock := new(mocks.UserMock)
	tokenMock := new(mocks.TokenMock)
	accessTTL := 24 * time.Hour
	refreshTTL := 168 * time.Hour
	tx := &transaction.Transactor{}
	email := "test@example.com"
	pass := "123qwe"
	_, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
	}
	svc := NewUserService(userMock, tokenMock, "test_secret", accessTTL, refreshTTL, tx)
	userMock.On("GetUserByEmail", mock.Anything, mock.Anything).Return(nil, errors.New("failed login"))

	_, _, err = svc.Login(context.Background(), email, pass)

	assert.Equal(t, domain.ErrInvalidCredentials, err)
	userMock.AssertExpectations(t)

}

func Test_LoginPassFail(t *testing.T) {
	userMock := new(mocks.UserMock)
	tokenMock := new(mocks.TokenMock)
	accessTTL := 24 * time.Hour
	refreshTTL := 168 * time.Hour
	tx := &transaction.Transactor{}
	email := "test@example.com"
	pass := "123qwe"
	fakePass := "qwe123"
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
	}
	svc := NewUserService(userMock, tokenMock, "test_secret", accessTTL, refreshTTL, tx)
	userMock.On("GetUserByEmail", mock.Anything, mock.Anything).Return(&domain.User{PasswordHash: string(hash)}, nil)

	access, refresh, err := svc.Login(context.Background(), email, fakePass)

	assert.Equal(t, domain.ErrInvalidCredentials, err)
	assert.Empty(t, access)
	assert.Empty(t, refresh)
	userMock.AssertExpectations(t)
	tokenMock.AssertNotCalled(t, "CreateToken", mock.Anything, mock.Anything)
}

func Test_GetUserByIdSuccess(t *testing.T) {
	userMock := new(mocks.UserMock)
	tokenMock := new(mocks.TokenMock)
	accessTTL := 24 * time.Hour
	refreshTTL := 168 * time.Hour
	tx := &transaction.Transactor{}
	svc := NewUserService(userMock, tokenMock, "test_secret", accessTTL, refreshTTL, tx)
	mockedUser := &domain.User{Id: "1", Email: "test@example.com"}
	userMock.On("GetUserById", mock.Anything, mock.Anything).Return(mockedUser, nil)

	user, err := svc.GetUserById(context.Background(), "1")

	assert.Equal(t, mockedUser, user)
	assert.NoError(t, err)
	userMock.AssertExpectations(t)
}

func Test_GetUserByIdFail(t *testing.T) {
	userMock := new(mocks.UserMock)
	tokenMock := new(mocks.TokenMock)
	accessTTL := 24 * time.Hour
	refreshTTL := 168 * time.Hour
	tx := &transaction.Transactor{}
	svc := NewUserService(userMock, tokenMock, "test_secret", accessTTL, refreshTTL, tx)
	userMock.On("GetUserById", mock.Anything, mock.Anything).Return(&domain.User{}, errors.New("failed to get user"))

	_, err := svc.GetUserById(context.Background(), "1")

	assert.Equal(t, errors.New("failed to get user"), err)
	userMock.AssertExpectations(t)
}

func Test_UpdateUserSuccess(t *testing.T) {
	userMock := new(mocks.UserMock)
	tokenMock := new(mocks.TokenMock)
	accessTTL := 24 * time.Hour
	refreshTTL := 168 * time.Hour
	tx := &transaction.Transactor{}
	model := &domain.User{}
	svc := NewUserService(userMock, tokenMock, "test_secret", accessTTL, refreshTTL, tx)
	userMock.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	err := svc.UpdateUser(context.Background(), model, "Ivan", "test@example.com")

	assert.Equal(t, nil, err)
	userMock.AssertExpectations(t)
}

func Test_UpdateUserFail(t *testing.T) {
	userMock := new(mocks.UserMock)
	tokenMock := new(mocks.TokenMock)
	accessTTL := 24 * time.Hour
	refreshTTL := 168 * time.Hour
	tx := &transaction.Transactor{}
	model := &domain.User{}
	svc := NewUserService(userMock, tokenMock, "test_secret", accessTTL, refreshTTL, tx)
	userMock.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("failed to update user"))

	err := svc.UpdateUser(context.Background(), model, "Ivan", "test@example.com")

	assert.Equal(t, errors.New("failed to update user"), err)
	userMock.AssertExpectations(t)
}

// Test_RefreshSuccess - успешное выполнение обновления токенов, mockUSer,
// mockToken - статичные заглушки с параметрами передаваемыми внутрь,
// sqlmock -используется, чтобы реальный код, мог выполняться внутри теста
// и не падать на nil(абстрактно говоря, мок транзакции),
// dbMock.ExpectBegin()/ExpectCommit() - регистрируют ожидание, что транзакция
// будет открыта и закоммичена (в нашем случае сценарий успешен), sqlx.NewDb
// оборачивает мокнутый драйвер в *sqlx.DB, dbMock.ExpectationsWereMet()
// проверяет, что эти ожидания реально выполнились и не осталось ошибок
func Test_RefreshSuccess(t *testing.T) {
	mockUser := new(mocks.UserMock)
	mockToken := new(mocks.TokenMock)
	accessTTL := 24 * time.Hour
	refreshTTL := 168 * time.Hour
	d := 1 * time.Hour
	db, dbMock, err := sqlmock.New()
	if err != nil {
		slog.Error("an error was not expected when opening a stub database connection", "error", err)
		return
	}
	defer func() {
		err := db.Close()
		if err != nil {
			slog.Error("failed to close db", "error", err)
		}
	}()
	dbMock.ExpectBegin()
	dbMock.ExpectCommit()
	tx := transaction.NewTransactor(sqlx.NewDb(db, "postgres"))
	svc := NewUserService(mockUser, mockToken, "test_secret", accessTTL, refreshTTL, tx)
	mockToken.On("GetTokenByHash", mock.Anything, mock.Anything).Return(&domain.RefreshToken{ExpiresAt: time.Now().Add(d)}, nil)
	mockUser.On("GetUserById", mock.Anything, mock.Anything).Return(&domain.User{}, nil)
	mockToken.On("RevokeToken", mock.Anything, mock.Anything).Return(nil)
	mockToken.On("CreateToken", mock.Anything, mock.Anything).Return(&domain.RefreshToken{}, nil)

	access, refresh, err := svc.Refresh(context.Background(), "test_token")

	assert.NotEmpty(t, access)
	assert.NotEmpty(t, refresh)
	assert.Equal(t, nil, err)
	assert.NoError(t, dbMock.ExpectationsWereMet())
	mockUser.AssertExpectations(t)
	mockToken.AssertExpectations(t)
}

func Test_RefreshGetHashFail(t *testing.T) {
	mockUser := new(mocks.UserMock)
	mockToken := new(mocks.TokenMock)
	accessTTL := 24 * time.Hour
	refreshTTL := 168 * time.Hour
	tx := &transaction.Transactor{}
	svc := NewUserService(mockUser, mockToken, "test:key", accessTTL, refreshTTL, tx)
	mockToken.On("GetTokenByHash", mock.Anything, mock.Anything).Return(&domain.RefreshToken{}, errors.New("failed to get hash"))

	_, _, err := svc.Refresh(context.Background(), "test_token")

	assert.Equal(t, errors.New("failed to get hash"), err)
	mockToken.AssertExpectations(t)
}

func Test_RefreshGetUserFail(t *testing.T) {
	mockUser := new(mocks.UserMock)
	mockToken := new(mocks.TokenMock)
	accessTTL := 24 * time.Hour
	refreshTTL := 168 * time.Hour
	d := 1 * time.Hour
	tx := &transaction.Transactor{}
	svc := NewUserService(mockUser, mockToken, "test", accessTTL, refreshTTL, tx)
	mockToken.On("GetTokenByHash", mock.Anything, mock.Anything).Return(&domain.RefreshToken{ExpiresAt: time.Now().Add(d)}, nil)
	mockUser.On("GetUserById", mock.Anything, mock.Anything).Return(&domain.User{}, errors.New("failed to get user"))

	_, _, err := svc.Refresh(context.Background(), "test_token")

	assert.Equal(t, errors.New("failed to get user"), err)
	mockToken.AssertExpectations(t)
	mockUser.AssertExpectations(t)
}

func Test_RefreshTokenRevoked(t *testing.T) {
	mockUser := new(mocks.UserMock)
	mockToken := new(mocks.TokenMock)
	accessTTL := 24 * time.Hour
	refreshTTL := 168 * time.Hour
	tx := &transaction.Transactor{}
	svc := NewUserService(mockUser, mockToken, "test_secret", accessTTL, refreshTTL, tx)
	revokedAt := time.Now()
	mockToken.On("GetTokenByHash", mock.Anything, mock.Anything).Return(&domain.RefreshToken{RevokedAt: &revokedAt}, nil)

	_, _, err := svc.Refresh(context.Background(), "test_token")

	assert.Equal(t, domain.ErrTokenRevoked, err)
	mockToken.AssertExpectations(t)
	mockUser.AssertNotCalled(t, "GetUserById", mock.Anything, mock.Anything)
}

func Test_RefreshTokenExpired(t *testing.T) {
	mockUser := new(mocks.UserMock)
	mockToken := new(mocks.TokenMock)
	accessTTL := 24 * time.Hour
	refreshTTL := 168 * time.Hour
	tx := &transaction.Transactor{}
	svc := NewUserService(mockUser, mockToken, "test_secret", accessTTL, refreshTTL, tx)
	mockToken.On("GetTokenByHash", mock.Anything, mock.Anything).Return(&domain.RefreshToken{ExpiresAt: time.Now().Add(-1 * time.Hour)}, nil)

	_, _, err := svc.Refresh(context.Background(), "test_token")

	assert.Equal(t, domain.ErrTokenExpired, err)
	mockToken.AssertExpectations(t)
	mockUser.AssertNotCalled(t, "GetUserById", mock.Anything, mock.Anything)
}

func Test_RefreshRevokeTokenFail(t *testing.T) {
	mockUser := new(mocks.UserMock)
	mockToken := new(mocks.TokenMock)
	accessTTL := 24 * time.Hour
	refreshTTL := 168 * time.Hour
	d := 1 * time.Hour
	db, dbMock, err := sqlmock.New()
	if err != nil {
		slog.Error("an error was not expected when opening a stub database connection", "error", err)
		return
	}
	defer func() {
		err := db.Close()
		if err != nil {
			slog.Error("failed to close db", "error", err)
		}
	}()
	dbMock.ExpectBegin()
	dbMock.ExpectRollback()
	tx := transaction.NewTransactor(sqlx.NewDb(db, "postgres"))
	svc := NewUserService(mockUser, mockToken, "test_secret", accessTTL, refreshTTL, tx)
	mockToken.On("GetTokenByHash", mock.Anything, mock.Anything).Return(&domain.RefreshToken{ExpiresAt: time.Now().Add(d)}, nil)
	mockUser.On("GetUserById", mock.Anything, mock.Anything).Return(&domain.User{}, nil)
	mockToken.On("RevokeToken", mock.Anything, mock.Anything).Return(errors.New("failed to revoke"))

	_, _, err = svc.Refresh(context.Background(), "test_token")

	assert.Equal(t, errors.New("failed to revoke"), err)
	assert.NoError(t, dbMock.ExpectationsWereMet())
	mockUser.AssertExpectations(t)
	mockToken.AssertExpectations(t)
	mockToken.AssertNotCalled(t, "CreateToken", mock.Anything, mock.Anything)
}

func Test_RefreshCreateTokenFail(t *testing.T) {
	mockUser := new(mocks.UserMock)
	mockToken := new(mocks.TokenMock)
	accessTTL := 24 * time.Hour
	refreshTTL := 168 * time.Hour
	d := 1 * time.Hour
	db, dbMock, err := sqlmock.New()
	if err != nil {
		slog.Error("an error was not expected when opening a stub database connection", "error", err)
		return
	}
	defer func() {
		err := db.Close()
		if err != nil {
			slog.Error("failed to close db", "error", err)
		}
	}()
	dbMock.ExpectBegin()
	dbMock.ExpectRollback()
	tx := transaction.NewTransactor(sqlx.NewDb(db, "postgres"))
	svc := NewUserService(mockUser, mockToken, "test_secret", accessTTL, refreshTTL, tx)
	mockToken.On("GetTokenByHash", mock.Anything, mock.Anything).Return(&domain.RefreshToken{ExpiresAt: time.Now().Add(d)}, nil)
	mockUser.On("GetUserById", mock.Anything, mock.Anything).Return(&domain.User{}, nil)
	mockToken.On("RevokeToken", mock.Anything, mock.Anything).Return(nil)
	mockToken.On("CreateToken", mock.Anything, mock.Anything).Return(nil, errors.New("failed to create token"))

	_, _, err = svc.Refresh(context.Background(), "test_token")

	assert.Equal(t, errors.New("failed to create token"), err)
	assert.NoError(t, dbMock.ExpectationsWereMet())
	mockUser.AssertExpectations(t)
	mockToken.AssertExpectations(t)
}

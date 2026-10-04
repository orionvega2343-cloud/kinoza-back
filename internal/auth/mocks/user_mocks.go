package mocks

import (
	"context"
	"kinoza-back/internal/auth/domain"

	"github.com/stretchr/testify/mock"
)

type UserMock struct {
	mock.Mock
}

func (m *UserMock) CreateUser(ctx context.Context, us *domain.User) (*domain.User, error) {
	args := m.Called(ctx, us)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *UserMock) GetUserById(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *UserMock) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *UserMock) UpdateUser(ctx context.Context, u *domain.User, name, email string) error {
	args := m.Called(ctx, u, name, email)
	return args.Error(0)
}

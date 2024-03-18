package mocks

import (
	"context"
	"github.com/MastersAndFans/masterfans-backend/pkg/models"
	"github.com/stretchr/testify/mock"
)

type UserRepoMock struct {
	mock.Mock
}

func (m *UserRepoMock) List(ctx context.Context, user *models.User) error {
	return nil
}

func (m *UserRepoMock) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepoMock) FindById(ctx context.Context, user *models.User) error {
	return nil
}

func (m *UserRepoMock) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

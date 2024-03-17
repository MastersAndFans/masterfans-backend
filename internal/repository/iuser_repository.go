package repository

import (
	"context"
	"github.com/MastersAndFans/masterfans-backend/pkg/models"
)

type IUserRepository interface {
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
}

package repository

import (
	"context"
	"github.com/MastersAndFans/masterfans-backend/pkg/models"
)

type IUserRepository interface {
	List(ctx context.Context) ([]models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindById(ctx context.Context, id int64) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
}

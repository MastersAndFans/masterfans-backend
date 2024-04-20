package repository

import (
	"context"
	"github.com/MastersAndFans/masterfans-backend/pkg/models"
)

type IAuctionRepository interface {
	List(ctx context.Context) ([]models.Auction, error)
	FindById(ctx context.Context, id int64) (*models.Auction, error)
	Create(ctx context.Context, auction *models.Auction) error
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, auction *models.Auction) (error)
}

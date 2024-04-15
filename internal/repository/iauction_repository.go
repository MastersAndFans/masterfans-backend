package repository

import (
	"context"
	"github.com/MastersAndFans/masterfans-backend/pkg/models"
)

type IAuctionRepository interface {
	List(ctx context.Context) ([]models.Auction, error)
	FindById(ctx context.Context, id int64) (*models.Auction, error)
	CreateAuction(ctx context.Context, auction *models.Auction) error
}

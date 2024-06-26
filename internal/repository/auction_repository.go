package repository

import (
	"context"
	"github.com/MastersAndFans/masterfans-backend/pkg/models"
	"gorm.io/gorm"
)

type AuctionRepository struct {
	db *gorm.DB
}

func NewAuctionRepository(db *gorm.DB) *AuctionRepository {
	return &AuctionRepository{db: db}
}

func (r *AuctionRepository) List(ctx context.Context) ([]models.Auction, error) {
	var auctions []models.Auction
	err := r.db.Preload("Proposer").Preload("Participants").Find(&auctions).Error
	if err != nil {
		return nil, err
	}
	return auctions, nil
}

func (r *AuctionRepository) FindById(ctx context.Context, id int64) (*models.Auction, error) {
	var auction models.Auction
	err := r.db.Preload("Proposer").Preload("Participants").WithContext(ctx).Where("id = ?", id).First(&auction).Error
	if err != nil {
		return nil, err
	}

	return &auction, nil
}

func (r *AuctionRepository) Create(ctx context.Context, auction *models.Auction) error {
	auctionRepo := NewRepository[models.Auction](r.db)

	return auctionRepo.Create(ctx, auction)
}

func (r *AuctionRepository) Delete(ctx context.Context, id int64) error {
	auction, err := r.FindById(ctx, id)
	if err != nil {
		return err
	}

	auctionRepo := NewRepository[models.Auction](r.db)
	return auctionRepo.Delete(ctx, auction)
}

func (r *AuctionRepository) Update(ctx context.Context, auction *models.Auction) (error) {
	err := r.db.Save(&auction).Error
	return err
}

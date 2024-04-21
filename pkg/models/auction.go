package models

import (
	"time"
)

type AuctionCategory int64

const (
	Undefined AuctionCategory = 0
	Carpenter AuctionCategory = 1
	Electrician AuctionCategory = 2
	Mechanic AuctionCategory = 3
	Mason AuctionCategory = 4
	InteriorDesigner AuctionCategory = 5
	SoftwareEngineer AuctionCategory = 6
	GraphicsDesigner AuctionCategory = 7
)

type Auction struct {
	ID            uint            `gorm:"primaryKey"          json:"auction_id"`
	ProposerID    uint            `gorm:"foreignKey"          json:"proposer_id"`
	Proposer      User            `                           json:"proposer"`
	Winner        *User           `gorm:"foreignKey:ID"       json:"winner"`
        Participants  []*User          `gorm:"Many2Many:auction_participants" json:"participants"`
	Active        bool            `                           json:"active"`
	StartingPrice int64           `                           json:"starting_price"`
	StartDate     time.Time       `                           json:"start_date"`
	EndDate       time.Time       `                           json:"end_date"`
	Title         string          `                           json:"title"`
	Description   string          `                           json:"description"`
	Category      AuctionCategory `                           json:"category"`
	City          string          `                           json:"city"`
}

package models

import "time"

type Auction struct {
	ID         uint `gorm:"primaryKey"          json:"auction_id"`
	ProposerID uint `gorm:"foreignKey"          json:"proposer_id"`
	Proposer   User `                           json:"proposer"`
	//WinnerID     uint
	//Winner       User      `gorm:"foreignKey:WinnerID" json:"winner"`
	Winner       User      `gorm:"foreignKey:ID" json:"winner"`
	Participants []User    `gorm:"foreignKey:ID"       json:"participants"`
	StartDate    time.Time `                           json:"start_date"`
	EndDate      time.Time `                           json:"end_date"`
	Title        string    `                           json:"title"`
	Description  string    `                           json:"description"`
}

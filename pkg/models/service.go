package models

type Service struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	MasterID    uint   `gorm:"not null"   json:"master_id"`
	Name        string `gorm:"not null"   json:"name"`
	Description string `gorm:"not null"   json:"description"`
	Price       uint   `gorm:"not null"   json:"price"`
}

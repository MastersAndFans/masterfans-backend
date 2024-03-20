package models

type Review struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	MasterID  uint   `gorm:"not null"   json:"master_id"`
	UserID    uint   `gorm:"not null"   json:"user_id"`
	Content   string `gorm:"not null"   json:"content"`
	Star      uint   `gorm:"not null"   json:"star"`
	CreatedAt string `                  json:"created_at"`
	UpdatedAt string `                  json:"updated_at"`
	DeletedAt string `                  json:"deleted_at"`
}

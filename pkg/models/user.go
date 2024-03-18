package models

import "time"

type User struct {
	ID          uint      `gorm:"primaryKey"         json:"user_id"`
	Email       string    `gorm:"unique;not null"    json:"email"`
	Password    string    `gorm:"not null"           json:"-"`
	Name        string    `gorm:"not null"           json:"name"`
	Surname     string    `gorm:"not null"           json:"surname"`
	BirthDate   time.Time `gorm:"type:date;not null" json:"birth_date"`
	PhoneNumber string    `                          json:"phone_number,omitempty"`
	IsMaster    bool      `gorm:"not null"           json:"is_master"`
	CreatedAt   time.Time `                          json:"created_at"`
	UpdatedAt   time.Time `                          json:"updated_at"`
}

package models

import "time"

type User struct {
	ID          uint      `gorm:"primaryKey"`
	Email       string    `gorm:"unique;not null"`
	Password    string    `gorm:"not null"`
	Name        string    `gorm:"not null"`
	Surname     string    `gorm:"not null"`
	BirthDate   time.Time `gorm:"type:date;not null"`
	PhoneNumber string
	IsMaster    bool `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

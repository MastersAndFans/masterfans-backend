package models

type Schedule struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	MasterID  uint   `gorm:"not null"   json:"master_id"`
	Day       string `gorm:"not null"   json:"day"`
	StartTime string `gorm:"not null"   json:"start_time"`
	EndTime   string `gorm:"not null"   json:"end_time"`
	IsBooked  bool   `gorm:"not null"   json:"is_booked"`
}

package models

type Schedule struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	MasterID   uint        `gorm:"not null"   json:"master_id"`
	Day        string      `gorm:"not null"   json:"day"`
	IsBooked   bool        `gorm:"not null"   json:"is_booked"`
	TimeRanges []TimeRange `gorm:"foreignKey:ScheduleID"`
}

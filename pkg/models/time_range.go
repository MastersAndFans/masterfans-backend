package models

import "time"

type TimeRange struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ScheduleID uint      `gorm:"not null"   json:"schedule_id"`
	StartTime  time.Time `gorm:"not null"   json:"start_time"`
	EndTime    time.Time `gorm:"not null"   json:"end_time"`
}

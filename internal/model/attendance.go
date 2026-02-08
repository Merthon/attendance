package model

import (
	"time"
)

type Attendance struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"userId"`
	Date      string    `gorm:"type:date;not null" json:"date"`
	CheckIn   *time.Time `json:"checkIn"`
	CheckOut  *time.Time `json:"checkOut"`
	Status    int       `gorm:"default:0" json:"status"` // 0-正常, 1-迟到, 2-早退
	Tags      string    `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	CheckInIP      string `json:"checkInIp"`
	CheckInDevice  string `json:"checkInDevice"`
	CheckOutIP     string `json:"checkOutIp"`
	CheckOutDevice string `json:"checkOutDevice"`

	// 关联 User
	User User `gorm:"foreignKey:UserID" json:"user"`
}
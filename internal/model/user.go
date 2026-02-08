package model

import "time"

// User 对应数据库 users 表
type User struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Username      string    `gorm:"unique;not null" json:"username"`
	Password      string    `gorm:"not null" json:"password,omitempty"` 
	RealName      string    `gorm:"not null" json:"realName"`
	Role          int       `gorm:"default:1" json:"role"`
	Department    string    `json:"department"`

	WorkStartTime string    `gorm:"default:'09:00:00'" json:"workStartTime"`
	WorkEndTime   string    `gorm:"default:'18:00:00'" json:"workEndTime"`

	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
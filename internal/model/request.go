package model

import "time"

// Request 对应数据库 requests 表
type Request struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null" json:"userId"`
	Type         int       `gorm:"not null" json:"type"` // 1-请假, 2-调休
	Category     string    `json:"category"`             // 事假/病假
	StartTime    time.Time `json:"startTime"`
	EndTime      time.Time `json:"endTime"`
	Reason       string    `json:"reason"`
	Status       int       `gorm:"default:0" json:"status"` // 0-待审批, 1-通过, 2-拒绝
	AdminComment string    `json:"adminComment"`
	CreatedAt    time.Time `json:"createdAt"`
    
    // 关联查询
	User         User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
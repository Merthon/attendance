package service

import (
	"attendance-system/internal/database"
	"attendance-system/internal/model"
	"errors"
)

// CreateRequest 提交申请
func CreateRequest(req *model.Request) error {
	// 简单校验
	if req.StartTime.After(req.EndTime) {
		return errors.New("结束时间不能早于开始时间")
	}
	
	// 入库
	return database.DB.Create(req).Error
}

// ApproveRequest 审批申请
func ApproveRequest(requestID uint, status int, comment string) error {
	var req model.Request
	if err := database.DB.First(&req, requestID).Error; err != nil {
		return errors.New("申请记录不存在")
	}

	// 更新状态
	req.Status = status
	req.AdminComment = comment
	
	return database.DB.Save(&req).Error
}

// GetAllRequests 获取所有申请
func GetAllRequests() ([]model.Request, error) {
	var requests []model.Request
	err := database.DB.Preload("User").Order("created_at desc").Find(&requests).Error
	return requests, err
}
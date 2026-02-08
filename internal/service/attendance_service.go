package service

import (
	"attendance-system/internal/database"
	"attendance-system/internal/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

// CheckIn 上班打卡 
func CheckIn(userID uint, ip string, device string) (*model.Attendance, error) {
	now := time.Now()
	today := now.Format("2006-01-02")

	// 1. 检查是否已打卡
	var exists model.Attendance
	err := database.DB.Where("user_id = ? AND date = ?", userID, today).First(&exists).Error
	if err == nil {
		return nil, errors.New("您今天已经打过上班卡了")
	}

	// 2. 获取用户信息 
	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 3. 动态判断迟到
	targetTimeStr := user.WorkStartTime
	if targetTimeStr == "" {
		targetTimeStr = "09:00:00" // 默认 9 点
	}
	// 解析今天的标准上班时间
	workStart, _ := time.ParseInLocation("2006-01-02 15:04:05", today+" "+targetTimeStr, time.Local)

	status := 0
	if now.After(workStart) {
		status = 1 // 迟到
	}

	// 4. 入库
	attendance := model.Attendance{
		UserID:        userID,
		Date:          today,
		CheckIn:       &now,
		Status:        status,
		CheckInIP:     ip,
		CheckInDevice: device,
	}
	attendance.User = user // 关联 User 信息以便返回

	if result := database.DB.Create(&attendance); result.Error != nil {
		return nil, result.Error
	}

	return &attendance, nil
}

// CheckOut 下班打卡
func CheckOut(userID uint, ip string, device string) (*model.Attendance, error) {
	now := time.Now()
	today := now.Format("2006-01-02")

	var attendance model.Attendance
	err := database.DB.Where("user_id = ? AND date = ?", userID, today).First(&attendance).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("请先打上班卡")
	}

	// 1. 获取用户信息
	var user model.User
	database.DB.First(&user, userID)

	// 2. 动态判断早退
	targetTimeStr := user.WorkEndTime
	if targetTimeStr == "" {
		targetTimeStr = "18:00:00" 
	}
	workEnd, _ := time.ParseInLocation("2006-01-02 15:04:05", today+" "+targetTimeStr, time.Local)

	// 只要还没到下班点，就算早退
	if now.Before(workEnd) {
		attendance.Status = 2 
	}

	// 3. 更新下班信息
	attendance.CheckOut = &now
	attendance.CheckOutIP = ip
	attendance.CheckOutDevice = device

	if result := database.DB.Save(&attendance); result.Error != nil {
		return nil, result.Error
	}

	attendance.User = user
	return &attendance, nil
}

// GetUserAttendance 获取指定用户的考勤
func GetUserAttendance(userID uint) ([]model.Attendance, error) {
	var records []model.Attendance
	err := database.DB.Where("user_id = ?", userID).Order("date desc").Find(&records).Error
	return records, err
}

// GetAllAttendances 获取所有考勤
func GetAllAttendances() ([]model.Attendance, error) {
	var records []model.Attendance
	// 预加载 User 信息
	err := database.DB.Preload("User").Order("date desc").Find(&records).Error
	return records, err
}
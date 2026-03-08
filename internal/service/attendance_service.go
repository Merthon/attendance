package service

import (
	"attendance-system/internal/database"
	"attendance-system/internal/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

// CheckIn 上班打卡 (按小时阵营判定)
func CheckIn(userID uint, ip string, device string) (*model.Attendance, error) {
	now := time.Now()
	today := now.Format("2006-01-02")

	// 1. 检查是否已打过上班卡
	var exists model.Attendance
	err := database.DB.Where("user_id = ? AND date = ?", userID, today).First(&exists).Error
	if err == nil {
		return nil, errors.New("您今天已经打过上班卡了")
	}

	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// ==========================================
	// 🛑 核心逻辑：按“小时”划分阵营，判定迟到
	// ==========================================
	h := now.Hour()
	m := now.Minute()
	timeCode := h*100 + m // 比如 9:16 变成 916

	status := 0 // 默认正常

	if h < 10 {
		// 【9点阵营】（包括8点多、9点多来的）
		if timeCode > 916 {
			status = 1 // 超过 9:16 算迟到
		}
	} else {
		// 【10点阵营】（10点及以后来的）
		if timeCode > 1016 {
			status = 1 // 超过 10:16 算迟到
		}
	}

	// 3. 入库
	attendance := model.Attendance{
		UserID:        userID,
		Date:          today,
		CheckIn:       &now,
		Status:        status,
		CheckInIP:     ip,
		CheckInDevice: device,
	}
	attendance.User = user

	if result := database.DB.Create(&attendance); result.Error != nil {
		return nil, result.Error
	}

	return &attendance, nil
}

// CheckOut 下班打卡 (按上班时的阵营决定下班时间)
func CheckOut(userID uint, ip string, device string) (*model.Attendance, error) {
	now := time.Now()
	today := now.Format("2006-01-02")

	// 1. 检查上班卡
	var attendance model.Attendance
	err := database.DB.Where("user_id = ? AND date = ?", userID, today).First(&attendance).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("请先打上班卡")
	}

	// 2. 检查下班卡防重复
	if attendance.CheckOut != nil {
		return nil, errors.New("您今天已经打过下班卡了，不可重复操作")
	}

	var user model.User
	database.DB.First(&user, userID)

	// ==========================================
	// 🛑 核心逻辑：读取早上的打卡时间，决定下班时间
	// ==========================================
	checkInTime := attendance.CheckIn
	inH := checkInTime.Hour()

	targetHour := 18 // 默认下班时间 18点

	// 如果早上是 10点及以后打的卡，下班时间变成 19点
	if inH >= 10 {
		targetHour = 19
	}

	targetEndTime := time.Date(now.Year(), now.Month(), now.Day(), targetHour, 0, 0, 0, now.Location())

	// 判断是否早退
	if now.Before(targetEndTime) {
		if attendance.Status == 1 {
			attendance.Status = 3 // 早上迟到，晚上又早退
		} else {
			attendance.Status = 2 // 正常上班，但早退
		}
	}

	// 更新下班信息
	attendance.CheckOut = &now
	attendance.CheckOutIP = ip
	attendance.CheckOutDevice = device

	if result := database.DB.Save(&attendance); result.Error != nil {
		return nil, result.Error
	}

	attendance.User = user
	return &attendance, nil
}

// GetUserAttendance 和 GetAllAttendances 
func GetUserAttendance(userID uint) ([]model.Attendance, error) {
	var records []model.Attendance
	err := database.DB.Where("user_id = ?", userID).Order("date desc").Find(&records).Error
	return records, err
}

func GetAllAttendances() ([]model.Attendance, error) {
	var records []model.Attendance
	err := database.DB.Preload("User").Order("date desc").Find(&records).Error
	return records, err
}
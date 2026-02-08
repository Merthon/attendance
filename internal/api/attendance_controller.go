package api

import (
	"attendance-system/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CheckInHandler 上班打卡
func CheckInHandler(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或Token无效"})
		return
	}

	// 获取信息
	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// 传给 Service 
	record, err := service.CheckIn(userId.(uint), clientIP, userAgent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "上班打卡成功！",
		"data":    record,
	})
}

// CheckOutHandler 下班打卡
func CheckOutHandler(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或Token无效"})
		return
	}

	// 自动获取信息
	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// 传给 Service (3个参数)
	record, err := service.CheckOut(userId.(uint), clientIP, userAgent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "下班打卡成功！",
		"data":    record,
	})
}

// GetMyAttendanceHandler 获取我的考勤
func GetMyAttendanceHandler(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	records, err := service.GetUserAttendance(userId.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": records})
}

// GetAllAttendancesHandler 获取全员考勤
func GetAllAttendancesHandler(c *gin.Context) {
	records, err := service.GetAllAttendances()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": records})
}
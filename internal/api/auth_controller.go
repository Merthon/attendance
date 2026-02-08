package api

import (
	"attendance-system/internal/model"
	"attendance-system/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginRequest 定义登录请求参数结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginHandler 用户登录
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	// 1. 绑定参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名或密码不能为空"})
		return
	}

	// 2. 调用 Service 层进行登录 
	token, user, err := service.Login(req.Username, req.Password)
	if err != nil {
		// 登录失败 (密码错误或用户不存在)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 3. 登录成功，返回 Token 和用户信息
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
		"user":    user,
	})
}

// RegisterHandler 用户注册
func RegisterHandler(c *gin.Context) {
	var user model.User
	// 1. 绑定参数
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := service.Register(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
		"data":    user, // 返回创建好的用户 (不含密码)
	})
}
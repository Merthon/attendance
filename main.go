package main

import (
	"attendance-system/internal/api"
	"attendance-system/internal/database"
	"attendance-system/internal/middleware" // 引入中间件
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 初始化数据库
	database.InitDB()

	// 2. 设置 Gin
	r := gin.Default()

	r.Static("/assets", "./dist/assets")

	// 公开路由组 (不需要登录就能访问)
	public := r.Group("/api/v1")
	{
		public.POST("/login", api.LoginHandler)
		public.POST("/register", api.RegisterHandler)
	}

	// 私有路由组 (必须登录带 Token 才能访问)
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware()) 
	{
		// 考勤相关
		protected.POST("/attendance/checkin", api.CheckInHandler)
		protected.POST("/attendance/checkout", api.CheckOutHandler)
		protected.GET("/attendance/my", api.GetMyAttendanceHandler)

		// 申请相关
		protected.POST("/request/create", api.CreateRequestHandler)
		protected.GET("/request/list", api.GetRequestsHandler)     // 管理员接口
		protected.POST("/request/approve", api.ApproveRequestHandler) // 管理员接口
		protected.GET("/admin/export", api.ExportExcelHandler)     // 管理员接口
		protected.GET("/admin/attendance", api.GetAllAttendancesHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		// 如果是 API 路径找不到，返回 JSON 404
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "接口不存在"})
			return
		}
		// 否则返回前端主页 (让 React 路由去处理)
		c.File("./dist/index.html")
	})

	// 启动
	port := "8080"
	fmt.Printf("系统已启动！\n")
	fmt.Printf("本机访问: http://localhost:%s\n", port)

	if err := r.Run(":" + port); err != nil {
		fmt.Println("启动失败:", err)
	}
}
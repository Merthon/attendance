package api

import (
	"attendance-system/internal/model"
	"attendance-system/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateRequestHandler 提交申请
func CreateRequestHandler(c *gin.Context) {
	var req model.Request
	// 1. 绑定前端传来的数据 (类型、时间、理由)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 2. 强制将 UserID 设为 Token 里的人
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	req.UserID = userId.(uint)

	// 3. 提交给 Service
	if err := service.CreateRequest(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "申请已提交"})
}

// GetRequestsHandler 获取所有申请
func GetRequestsHandler(c *gin.Context) {
	requests, err := service.GetAllRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": requests})
}

// ApproveRequest 定义审批时的参数结构
type ApproveRequest struct {
	RequestID uint   `json:"requestId"`
	Status    int    `json:"status"` // 1-通过, 2-拒绝
	Comment   string `json:"comment"`
}

// ApproveRequestHandler 审批申请
func ApproveRequestHandler(c *gin.Context) {
	var req ApproveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := service.ApproveRequest(req.RequestID, req.Status, req.Comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "审批完成"})
}

// ExportExcelHandler 导出 Excel 
func ExportExcelHandler(c *gin.Context) {
	f, err := service.GenerateExcel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成失败"})
		return
	}

	// 设置响应头
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=attendance_report.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")

	// 将文件流写入响应
	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "下载失败"})
	}
}
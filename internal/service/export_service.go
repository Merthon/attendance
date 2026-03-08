package service

import (
	"attendance-system/internal/database"
	"attendance-system/internal/model"
	"github.com/xuri/excelize/v2"
	"strconv"
)

// GenerateExcel 
func GenerateExcel() (*excelize.File, error) {
    f := excelize.NewFile()
    
    // 1. 创建表头
    headers := []string{"ID", "员工姓名", "日期", "上班时间", "下班时间", "状态"}
    for i, h := range headers {
        // A1, B1, C1...
        cell := string(rune('A'+i)) + "1"
        f.SetCellValue("Sheet1", cell, h)
    }

    // 2. 查询数据 (🌟 优化：按日期倒序排，最新的考勤在最上面)
    var records []model.Attendance
    database.DB.Preload("User").Order("date desc").Find(&records) 

    // 3. 填充数据
    for i, r := range records {
        row := strconv.Itoa(i + 2) 
        f.SetCellValue("Sheet1", "A"+row, r.ID)
        
        // 🌟 修复1：填充真实的员工姓名，如果没填真名就用登录名兜底
        empName := r.User.RealName
        if empName == "" {
            empName = r.User.Username
        }
        f.SetCellValue("Sheet1", "B"+row, empName) 
        
        f.SetCellValue("Sheet1", "C"+row, r.Date)
        
        checkIn := ""
        if r.CheckIn != nil {
            checkIn = r.CheckIn.Format("15:04:05")
        }
        f.SetCellValue("Sheet1", "D"+row, checkIn)

        checkOut := ""
        if r.CheckOut != nil {
            checkOut = r.CheckOut.Format("15:04:05")
        }
        f.SetCellValue("Sheet1", "E"+row, checkOut)
        
        // 🌟 修复2：补全所有的状态判定
        statusText := "正常"
        if r.Status == 1 { statusText = "迟到" }
        if r.Status == 2 { statusText = "早退" }
        if r.Status == 3 { statusText = "迟到且早退" } // 👈 加上我们之前设计的双重惩罚状态
        
        f.SetCellValue("Sheet1", "F"+row, statusText)
    }

    return f, nil
}
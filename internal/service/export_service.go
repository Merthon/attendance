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

	// 2. 查询数据]
	var records []model.Attendance
	database.DB.Preload("User").Find(&records) 

	// 3. 填充数据
	for i, r := range records {
		row := strconv.Itoa(i + 2) 
		f.SetCellValue("Sheet1", "A"+row, r.ID)
		f.SetCellValue("Sheet1", "B"+row, r.UserID) 
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
		
		statusText := "正常"
		if r.Status == 1 { statusText = "迟到" }
		if r.Status == 2 { statusText = "早退" }
		f.SetCellValue("Sheet1", "F"+row, statusText)
	}

	return f, nil
}
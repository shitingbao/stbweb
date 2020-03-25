package excel

import (
	"fmt"
	"stbweb/core"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/tealeg/xlsx"
)

//getExcelAllCell 使用tealeg解析，时间不完美
func getExcelAllCell(fileURL string) error {
	var mySlice [][][]string
	mySlice, err := xlsx.FileToSlice(fileURL)
	if err != nil {
		return err
	}
	for _, v := range mySlice {
		for _, val := range v {
			for _, value := range val {
				core.LOG.Info("vel:", value)
			}
		}
	}
	return nil
}

//getExcelRows 使用360解析，时间不完美
func getExcelRows(excelURL string) error {
	xlsx, err := excelize.OpenFile(excelURL)
	if err != nil {
		return err
	}

	rows, err := xlsx.GetRows("Sheet1")
	if err != nil {
		return err
	}
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
	}
	return nil
}

// Export 导出
func Export() {
	getExcelRows("./file/export.xlsx")
}

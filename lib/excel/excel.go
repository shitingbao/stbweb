package excel

import (
	"fmt"
	"stbweb/core"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/Sirupsen/logrus"
	"github.com/tealeg/xlsx"
)

//Excel 数据对象
type Excel struct {
	FileName   string //文件名称，不需要文件后缀
	SheetNum   int    //sheet的数量,与sheet的数据数组长度对应
	SheetDatas []SheetData
}

//SheetData sheet数据内容
type SheetData struct {
	Title []string            //标题列
	Rows  []map[string]string //行，每一行的数据以标题为标准key，存储
}

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

//createExcel 新建一个excel
func (e *Excel) createExcel() {
	if len(e.SheetDatas) != e.SheetNum {
		return
	}
	file := xlsx.NewFile()

	for i, v := range e.SheetDatas {
		sheetName := "Sheet" + strconv.Itoa(i)
		// Create a new sheet.
		sheet, err := file.AddSheet(sheetName)
		if err != nil {
			return
		}
		for idv, dv := range v.Rows {
			row := sheet.AddRow()
			row.SetHeightCM(1) //设置每行的高度
			for key, val := range dv {
				cell := row.AddCell()
				if idv == 0 {
					cell.Value = key
				} else {
					cell.Value = val
				}
			}
		}
	}

	if err := file.Save(e.FileName + ".xlsx"); err != nil {
		core.LOG.WithFields(logrus.Fields{"excel": err}).Error("file")
	}
}

// Export 导出
func Export() {
	getExcelRows("./file/export.xlsx")
}

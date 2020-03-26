package excel

import (
	"fmt"
	"stbweb/core"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/Sirupsen/logrus"
	"github.com/tealeg/xlsx"
)

//excel 数据对象
type excel struct {
	FileName   string      //文件名称，不需要文件后缀
	SheetDatas []sheetData //sheet的数量,与sheet的数据数组长度对应
}

//sheetData sheet数据内容
type sheetData struct {
	Rows []map[string]string //行，每一行的数据以标题为标准key，存储
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
func (e *excel) createExcel() {
	file := xlsx.NewFile()
	for i, v := range e.SheetDatas {
		sheetName := "Sheet" + strconv.Itoa(i+1)
		// Create a new sheet.
		sheet, err := file.AddSheet(sheetName)
		if err != nil {
			return
		}
		for idv, dv := range v.Rows {
			var trow *xlsx.Row
			if idv == 0 { //第一行数据需要另外新增一行作为标题
				trow = sheet.AddRow()
				trow.SetHeightCM(1) //设置每行的高度
			}

			row := sheet.AddRow()
			row.SetHeightCM(1)         //设置每行的高度
			for key, val := range dv { //这里必须在一个for内部，不然key和val对应不上
				if idv == 0 {
					tcell := trow.AddCell()
					tcell.Value = key
				}
				cell := row.AddCell()
				cell.Value = val
			}
		}
	}

	if err := file.Save(e.FileName + ".xlsx"); err != nil {
		core.LOG.WithFields(logrus.Fields{"excel": err}).Error("file")
	}
}

//CreateExcel 创建一个excel
//name为文件名，不需要后缀
//后续参数为每个页面的数据，每一个rowData参数对应一个sheet页面,内部数据就是sheet的数据
//列名称就是map的key值
//example: CreateExcel("example",data1,data2)
//执行后生成文件名称为example.xlsx,内部有两个sheet页，sheet1数据内容为dat1，sheet2数据内容为dat2
func CreateExcel(name string, rowData ...[]map[string]string) {
	sheetDatas := []sheetData{}
	for _, v := range rowData {
		sd := sheetData{}
		sd.Rows = v
		sheetDatas = append(sheetDatas, sd)
	}
	el := excel{
		FileName:   name,
		SheetDatas: sheetDatas,
	}
	el.createExcel()
}

// Export 导出
func Export() {
	getExcelRows("./file/export.xlsx")
}

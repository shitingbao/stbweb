package excel

import "github.com/tealeg/xlsx"

func getExcelAllCell(fileURL string) {
	var mySlice [][][]string
	var value string
	mySlice = xlsx.FileToSlice(fileURL)
	value = mySlice[0][0][0]
}

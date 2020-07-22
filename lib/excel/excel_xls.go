package excel

import "github.com/extrame/xls"

//xls文件解析，反馈数据二维数组
func openXlsFile(url string) [][]string {
	dataList := [][]string{}
	xlFile, err := xls.Open(url, "utf-8")
	if err != nil {
		panic(err)
	}
	if sheet1 := xlFile.GetSheet(0); sheet1 != nil {
		for idx := 0; idx <= int(sheet1.MaxRow); idx++ {
			row := sheet1.Row(idx)
			data := []string{}
			for i := 0; i < row.LastCol(); i++ {
				data = append(data, row.Col(i))
			}
			dataList = append(dataList, data)
		}
	}
	return dataList
}

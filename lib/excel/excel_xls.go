package excel

import (
	"io/ioutil"
	"path"
	"strings"

	"github.com/extrame/xls"
)

//xls文件解析，反馈数据二维数组,获取的数据内容
//有一定bug，在单元格在超过一定行数后，会自动分行（一个单元格中信息突然变成两列），后面一个单元格内容无法获取到，越到后面丢失的数据越多
//这bug难以处理，就把xls转成xlsx在进行操作，方法待定
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

var baseURL = "C:/Users/87125/Desktop/report/2020-09-09app_data"

//输入保存路径和文件前缀
func load(url, pro string) {
	urlList, err := getFileList(baseURL)
	if err != nil {
		panic(err)
	}
	for _, v := range urlList {
		fileName := strings.Split(v, ".")
		var data [][]string
		switch fileName[1] {
		case "xls":
			data = openXlsFile(path.Join(baseURL, v))
		}
		fName := strings.Replace(fileName[0], " ", "-", -1)
		fName = strings.Replace(fileName[0], "_", "-", -1)
		if err := createCsvFile(path.Join(url, pro+fName+".csv"), "", data); err != nil {
			panic(err)
		}
	}
}

//返回所有文件地址列表
func getFileList(dir string) ([]string, error) {
	fList, err := ioutil.ReadDir(dir)
	if err != nil {
		return []string{}, err
	}
	fileList := []string{}
	for _, v := range fList {
		fileList = append(fileList, v.Name())
		// log.Println(v.Name())
	}
	return fileList, nil
}

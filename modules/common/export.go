package common

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"stbweb/core"
	"stbweb/lib/excel"
	"strings"
)

type export struct{}

func init() {
	core.RegisterFun("export", new(export), false)
}

//Get 业务处理,get请求的例子
func (ap *export) Get(arge *core.ElementHandleArgs) {
	if arge.APIInterceptionGet("excel", nil, excelExport) ||
		arge.APIInterceptionGet("excelparse", nil, excelParse) ||
		arge.APIInterceptionGet("xlsparse", nil, xlsParse) ||
		arge.APIInterceptionGet("csvparse", nil, csvParase) {
		return
	}
}

func (ap *export) Post(p *core.ElementHandleArgs) {
	p.Req.ParseMultipartForm(20 << 20)
	if p.Req.MultipartForm == nil {
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": false, "msg": "MultipartForm is null"})
		return
	}
	_, fHeader, err := p.Req.FormFile("file")
	if err != nil {
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": false, "msg": "get file have error"})
		return
	}
	fileAdree, err := getUpdateFile(fHeader)
	if err != nil {
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": false, "msg": err.Error()})
		return
	}
	sep, createSep, createFileType, isGBK, isCreateGBK := getFormValues(p)
	resURL := ""
	switch createFileType {
	case "csv":
		fileAdree, err := fileToCsv(fileAdree, sep, createSep, isGBK, isCreateGBK)
		if err != nil {
			core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": false, "msg": err.Error()})
			return
		}
		resURL = fileAdree
	case "excel":
		fileAdree, err := fileToExcel(fileAdree, sep, isGBK)
		if err != nil {
			core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": false, "msg": err.Error()})
			return
		}
		resURL = fileAdree
	}
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": true, "url": resURL})
}

//sep文件分割符，传入文件gbk是否gbk格式，createSep生成文件的分割符， isCreateGBK生成文件是否gbk格式，createFileType 内容为csv或者excel，代表生成哪种文件
func getFormValues(p *core.ElementHandleArgs) (sep, createSep, createFileType string, isGBK, isCreateGBK bool) {
	for k, v := range p.Req.MultipartForm.Value { //获取表单字段
		switch k {
		case "sep":
			sep = v[0]
		case "gbk":
			if v[0] == "true" {
				isGBK = true
			}
		case "createSep":
			createSep = v[0]
		case "isCreateGBK":
			if v[0] == "true" {
				isCreateGBK = true
			}
		case "createFileType":
			createFileType = v[0]
		}
	}
	return
}

//isGBK true标识使用gbk解析,isCreateGBK标识生成的csv是否用gbk，true代表使用,createSep标识生成文件的间隔符
//只能解析xlsx , csv , txt三种文件，都生成csv
func fileToCsv(fileURL, sep, createSep string, isGBK, isCreateGBK bool) (string, error) {
	fileData := [][]string{}
	switch path.Ext(fileURL) {
	case ".xlsx":
		fd, err := excel.ExportParse(fileURL)
		if err != nil {
			return "", err
		}
		fileData = fd
	case ".csv", ".txt":
		fileData = excel.PaseCscOrTxt(fileURL, sep, isGBK)
	default:
		return "", errors.New("file type error")
	}
	fileName := strings.TrimSuffix(path.Base(fileURL), path.Ext(fileURL))
	fileAdree := path.Join(core.DefaultFilePath, fileName+".csv")
	switch {
	case isCreateGBK:
		if err := excel.CreateGBKCsvFile(fileAdree, createSep, fileData); err != nil {
			return "", err
		}
	default:
		if err := excel.CreateCsvFile(fileAdree, createSep, fileData); err != nil {
			return "", err
		}
	}
	return fileAdree, nil
}

func fileToExcel(fileURL, sep string, isGBK bool) (string, error) {
	fileData := excel.PaseCscOrTxt(fileURL, sep, isGBK)
	fileName := strings.TrimSuffix(path.Base(fileURL), path.Ext(fileURL))
	fileAdree := path.Join(core.DefaultFilePath, fileName+".xlsx")
	if err := excel.CreateExcelUseList(fileAdree, fileData); err != nil {
		return "", err
	}
	return fileAdree, nil
}

func excelExport(pa interface{}, content *core.ElementHandleArgs) error {
	rowData := []map[string]string{}
	da := make(map[string]string)
	da["one"] = "one"
	da["Two"] = "two"
	da["三"] = "三"
	da["4"] = "4"
	da["date"] = "1994-08-01"
	rowData = append(rowData, da)
	rowDatat := []map[string]string{}
	dc := make(map[string]string)
	dc["asdf"] = "asdf"
	dc["asdf"] = "asdf"
	dc["三"] = "三"
	dc["4"] = "4"
	dc["date"] = "1994-08-01"
	rowDatat = append(rowDatat, dc)
	if err := excel.CreateExcel("stb", rowData, rowDatat); err != nil {
		core.SendJSON(content.Res, http.StatusOK, core.SendMap{"msg": err.Error()})
		return err
	}
	core.SendJSON(content.Res, http.StatusOK, core.SendMap{"msg": "this is excel get"})
	return nil
}

func excelParse(pa interface{}, content *core.ElementHandleArgs) error {
	res, err := excel.ExportParse("./assets/aa.xls", "Sheet1")
	if err != nil {
		return err
	}
	core.SendJSON(content.Res, http.StatusOK, res)
	return nil
}

func xlsParse(pa interface{}, content *core.ElementHandleArgs) error {
	res, err := excel.ExportParse("./assets/stb.xlsx", "Sheet2")
	if err != nil {
		return err
	}
	core.SendJSON(content.Res, http.StatusOK, res)
	return nil
}

func csvParase(pa interface{}, content *core.ElementHandleArgs) error {
	reMapE, reMapQ, err := excel.ComparisonFile("./assets/gg.csv", "./assets/gs.csv")
	if err != nil {
		core.SendJSON(content.Res, http.StatusOK, core.SendMap{"err": err.Error()})
		return nil
	}

	core.SendJSON(content.Res, http.StatusOK, core.SendMap{"E": *reMapE, "Q": *reMapQ})
	return nil
}

//获取表单中的文件，保存至默认路径并反馈保存的文件名
func getUpdateFile(file *multipart.FileHeader) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()
	if err := os.MkdirAll(core.DefaultFilePath, os.ModePerm); err != nil {
		return "", err
	}
	fileAdree := path.Join(core.DefaultFilePath, file.Filename)
	fl, err := os.Create(fileAdree)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(fl, f); err != nil {
		return "", err
	}
	return fileAdree, nil
}

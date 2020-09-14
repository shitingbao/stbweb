package common

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"stbweb/core"
	"stbweb/lib/excel"
	"stbweb/lib/formopera"
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
	fileHeaderList := formopera.GetAllFormFiles(p.Req)
	resFilePaths := []string{}
	for _, v := range fileHeaderList {
		path, err := getCsvFile(v)
		if err != nil {
			core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": false, "msg": err.Error()})
			return
		}
		resFilePaths = append(resFilePaths, path)
	}
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": true, "filepaths": resFilePaths})
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

//获取表单中的文件，保存至默认路径并反馈保存的文件路径
func getCsvFile(file *multipart.FileHeader) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()
	// ft := path.Ext(file.Filename assets)
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
	str, _ := os.Executable() //使用绝对路径给前端
	workDir := filepath.Dir(str)

	return strings.Replace(path.Join(workDir, fileAdree), "\\", "/", -1), nil
}

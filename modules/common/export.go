package common

import (
	"net/http"
	"stbweb/core"
	"stbweb/lib/excel"
)

type export struct{}

func init() {
	core.RegisterFun("export", new(export), false)
}

//Get 业务处理,get请求的例子
func (ap *export) Get(arge *core.ElementHandleArgs) {
	if arge.APIInterceptionGet("excel", nil, excelExport) ||
		arge.APIInterceptionGet("excelparse", nil, excelparse) ||
		arge.APIInterceptionGet("csvparse", nil, csvParase) {
		return
	}
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

func excelparse(pa interface{}, content *core.ElementHandleArgs) error {
	res, err := excel.ExportParse("./assets/stb.xlsx", "Sheet2")
	if err != nil {
		return err
	}
	core.SendJSON(content.Res, http.StatusOK, res)
	return nil
}

//每列的对象
type lod struct {
	Class1 string
	Class2 string
	Name   string
	Data   string
	Unit   string
}

type recode struct {
	Name string
	Mark int
}

func csvParase(pa interface{}, content *core.ElementHandleArgs) error {

	// recode := []string{}
	recodeE := []int{}
	recodeQ := []int{}
	recordE, err := excel.LoadCsvCfg("./assets/gg.csv")
	if err != nil {
		return err
	}

	reMapE := getRes(recordE)

	recordQ, err := excel.LoadCsvCfg("./assets/gs.csv")
	if err != nil {
		return err
	}

	reMapQ := getRes(recordQ)

	for k, v := range *reMapE {
		for _, vel := range *reMapQ {
			if v == vel {
				// EM := strconv.Itoa(k)
				// QM := strconv.Itoa(key)
				// recode = append(recode, "E:"+EM+"==Q:"+QM)
				recodeE = append(recodeE, k)
				break
			}
		}
	}

	for k, v := range *reMapQ {
		for _, vel := range *reMapE {
			if v == vel {
				// EM := strconv.Itoa(k)
				// QM := strconv.Itoa(key)
				// recode = append(recode, "E:"+EM+"==Q:"+QM)
				recodeQ = append(recodeQ, k)
				break
			}
		}
	}

	for _, v := range recodeE {
		delete(*reMapE, v)
	}

	for _, v := range recodeQ {
		delete(*reMapQ, v)
	}
	core.SendJSON(content.Res, http.StatusOK, core.SendMap{"E": *reMapE, "Q": *reMapQ})
	return nil
}

func getRes(recordE [][]string) *map[int]lod {
	reMapE := make(map[int]lod) //int记录行数
	for i, v := range recordE {
		da := lod{}
		for idx, vel := range v {
			switch idx {
			case 0:
				da.Class1 = vel
			case 1:
				da.Class2 = vel
			case 2:
				da.Name = vel
			case 3:
				da.Data = vel
			case 4:
				da.Unit = vel
			}
		}
		reMapE[i] = da
	}
	return &reMapE
}

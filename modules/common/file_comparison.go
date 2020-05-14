package common

import (
	"io"
	"net/http"
	"os"
	"path"
	"stbweb/core"
	"stbweb/lib/comparison"
	"stbweb/lib/formopera"

	"github.com/Sirupsen/logrus"
	"github.com/pborman/uuid"
)

type fileComparison struct{}

func init() {
	core.RegisterFun("filecomparison", new(fileComparison), false)
}

type comparisonParam struct {
	CompareFile comparison.ParisonFileObject
	AimFile     comparison.ParisonFileObject
}

//post中分api请求比对和表单比对
// api中直接输入文件路径
// 表单中获取文件以及相关文件标识
// 标识参数left，lft，lsep，listitle / right，rft，rsep，ristitle
//分别是两个文件的相关标识（左右）：文件，文件类型，文件分隔标识符，是否是标题
func (f *fileComparison) Post(p *core.ElementHandleArgs) {
	if p.APIInterceptionPost("comparison", new(comparisonParam), comparisonFileCommon) {
		return
	}
	rPath, err := getFormFile("lft", "left", "lsep", "listitle", p)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Err": err.Error()}).Error("getFormFile")
	}
	lPath, err := getFormFile("rft", "right", "rsep", "ristitle", p)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Err": err.Error()}).Error("getFormFile")
	}

	res, err := comparison.FileComparison(rPath, lPath)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Err": err.Error()}).Error("FileComparison")
		return
	}
	core.SendJSON(p.Res, http.StatusOK, res)
}

func comparisonFileCommon(param interface{}, p *core.ElementHandleArgs) error {
	pa := param.(*comparisonParam)
	res, err := comparison.FileComparison(pa.CompareFile, pa.AimFile)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Err": err.Error()}).Error("FileComparison")
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, res)
	return nil
}

//根据传入文件名称标识，文件类型标识，从formdata中获取文件
func getFormFile(tkey, fileName, sep, isTitle string, p *core.ElementHandleArgs) (comparison.ParisonFileObject, error) {
	lFileType := ""
	fSep := ""
	fIsTitle := false
	p.Req.ParseMultipartForm(20 << 20)
	for k, v := range p.Req.MultipartForm.Value { //获取表单字段
		if k == tkey && len(v) > 0 {
			lFileType = v[0]
		}
		if k == sep && len(v) > 0 {
			fSep = v[0]
		}
		if k == isTitle && len(v) > 0 && v[0] == "true" {
			fIsTitle = true
		}
	}
	lfile, err := formopera.GetFormOnceFile(fileName, p.Req)
	if err != nil {
		return comparison.ParisonFileObject{}, err
	}
	defer lfile.Close()
	if err := os.MkdirAll(core.DefaultFilePath, os.ModePerm); err != nil {
		return comparison.ParisonFileObject{}, err
	}
	fileAdree := path.Join(core.DefaultFilePath, uuid.NewUUID().String()+lFileType)
	fl, err := os.Create(fileAdree)
	if err != nil {
		return comparison.ParisonFileObject{}, err
	}
	if _, err := io.Copy(fl, lfile); err != nil {
		return comparison.ParisonFileObject{}, err
	}
	fInfo := comparison.ParisonFileObject{
		FileName: fileAdree,
		Sep:      fSep,
		IsTitle:  fIsTitle,
	}

	return fInfo, nil
}

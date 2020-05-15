package common

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"stbweb/core"
	"stbweb/lib/comparison"

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
	rPath, lPath, err := getFormFile(p)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Err": err.Error()}).Error("getFormFile")
		return
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
func getFormFile(p *core.ElementHandleArgs) (comparison.ParisonFileObject, comparison.ParisonFileObject, error) {
	p.Req.ParseMultipartForm(20 << 20)
	leftObject, rightObject := comparison.ParisonFileObject{}, comparison.ParisonFileObject{}
	if p.Req.MultipartForm == nil {
		return leftObject, rightObject, errors.New("form is nil")
	}
	for k, v := range p.Req.MultipartForm.Value { //获取表单字段
		switch k {
		case "lsep":
			leftObject.Sep = v[0]
		case "listitle":
			leftObject.IsTitle = false
		case "rsep":
			rightObject.Sep = v[0]
		case "ristitle":
			rightObject.IsTitle = false
		}
	}
	ladree, err := getSaveFilePath("left", p)
	if err != nil {
		return leftObject, rightObject, err
	}
	leftObject.FileName = ladree
	radree, err := getSaveFilePath("right", p)
	if err != nil {
		return leftObject, rightObject, err
	}
	rightObject.FileName = radree
	return leftObject, rightObject, nil
}

//获取表单中的文件，保存至默认路径并反馈保存的文件路径
func getSaveFilePath(fileName string, p *core.ElementHandleArgs) (string, error) {
	_, file, err := p.Req.FormFile(fileName)
	if err != nil {
		return "", err
	}
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	ft := path.Ext(file.Filename)
	if err := os.MkdirAll(core.DefaultFilePath, os.ModePerm); err != nil {
		return "", err
	}
	fileAdree := path.Join(core.DefaultFilePath, uuid.NewUUID().String()+ft)
	fl, err := os.Create(fileAdree)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(fl, f); err != nil {
		return "", err
	}
	return fileAdree, nil
}

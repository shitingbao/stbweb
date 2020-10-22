package common

import (
	"net/http"
	"stbweb/core"
	"stbweb/lib/formopera"
	"stbweb/lib/images"
	imagetowordapi "stbweb/lib/imagetowordAPI"
	"time"

	"github.com/sirupsen/logrus"
)

//ImageWord 业务类
type ImageWord struct{}

func init() {
	core.RegisterFun("image", new(ImageWord), false)
}

//Post 图片转文字
func (im *ImageWord) Post(p *core.ElementHandleArgs) {
	logrus.Info("image to word API")
	imageURLs, err := getFileHands(p)
	if err != nil {
		core.SendJSON(p.Res, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := imagesOpera(imageURLs, p)
	if err != nil {
		core.SendJSON(p.Res, http.StatusOK, err.Error())
		return
	}

	core.SendJSON(p.Res, http.StatusOK, result)
}

//获取表单内图片保存，并反馈对应所有图片路径
func getFileHands(p *core.ElementHandleArgs) ([]string, error) {
	imageURLs := []string{}
	fileHands := formopera.GetAllFormFiles(p.Req)
	for _, v := range fileHands {
		file, err := v.Open()
		if err != nil {
			return imageURLs, err
		}
		imageURL, err := images.ByteToImage(core.DefaultFilePath, file)
		if err != nil {
			return imageURLs, err
		}
		imageURLs = append(imageURLs, imageURL)
		file.Close()
	}
	return imageURLs, nil
}

//imagesOpera 传入图片路径，亲求三方接口反馈文字对象,需要先检查token可用性
func imagesOpera(imageURLs []string, p *core.ElementHandleArgs) ([]imagetowordapi.AcceptResultWord, error) {
	result := []imagetowordapi.AcceptResultWord{}
	token, err := imagetowordapi.CheckTokenEffect(core.WebConfig.AccessTokenDate)
	if err != nil {
		return result, err
	}
	if token != "" {
		core.WebConfig.AccessToken = token
		core.WebConfig.AccessTokenDate = time.Now().Format("2006-01-02 15:04:05")
		core.WebConfig.SaveConfig()
	}

	for _, v := range imageURLs {
		imagesBase64 := []string{}
		base64, err := images.ImageToBase64(v)
		if err != nil {
			return result, err
		}
		imagesBase64 = append(imagesBase64, base64)
		res, err := imagetowordapi.GetImageWord(core.WebConfig.BaidubceAddress, core.WebConfig.AccessToken, core.WebConfig.AccessTokenDate, imagesBase64)
		if err != nil {
			return result, err
		}
		result = append(result, res)
	}
	return result, nil
}

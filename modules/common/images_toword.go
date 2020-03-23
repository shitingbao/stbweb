package common

import (
	"net/http"
	"stbweb/core"
	"stbweb/lib/formopera"
	"stbweb/lib/images"
	imagetowordapi "stbweb/lib/imagetowordAPI"
)

//ImageWord 业务类
type ImageWord struct{}

type image struct {
	Image []string
}

func init() {
	core.RegisterFun("image", new(ImageWord))
}

//Post 图片转文字
func (im *ImageWord) Post(p *core.ElementHandleArgs) {
	imageURLs := []string{}
	fileHands := formopera.GetAllFormFiles(p.Req)
	for _, v := range fileHands {
		file, err := v.Open()
		if err != nil {
			core.SendJSON(p.Res, http.StatusInternalServerError, err.Error())
			return
		}
		imageURL, err := images.ByteToImage(file)
		if err != nil {
			core.SendJSON(p.Res, http.StatusInternalServerError, err.Error())
			return
		}
		imageURLs = append(imageURLs, imageURL)
		file.Close()
	}
	imagesBase64 := []string{}
	for _, v := range imageURLs {
		base64, err := images.ImageToBase64(v)
		if err != nil {
			core.SendJSON(p.Res, http.StatusOK, core.SendMap{"err": err.Error()})
			return
		}
		imagesBase64 = append(imagesBase64, base64)
	}
	core.LOG.Info(len(imagesBase64))
	res, err := imagetowordapi.GetImageWord(imagesBase64)
	if err != nil {
		core.SendJSON(p.Res, http.StatusOK, err.Error())
		return
	}
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"data": res})
}

func imagesOpera(pa interface{}, p *core.ElementHandleArgs) error {
	param := pa.(*image)
	core.SendJSON(p.Res, http.StatusOK, param)
	return nil
}

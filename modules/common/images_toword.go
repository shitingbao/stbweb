package common

import (
	"net/http"
	"stbweb/core"
	"stbweb/lib/formopera"
	"stbweb/lib/images"
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
		core.LOG.Info(imageURL)
		file.Close()
	}
	core.SendJSON(p.Res, http.StatusOK, "success")
}

func imagesOpera(pa interface{}, p *core.ElementHandleArgs) error {
	param := pa.(*image)
	core.SendJSON(p.Res, http.StatusOK, param)
	return nil
}

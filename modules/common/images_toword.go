package common

import (
	"net/http"
	"stbweb/core"
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
	if p.APIInterceptionPost("imageword", new(image), imagesOpera) {
		return
	}
}

func imagesOpera(pa interface{}, p *core.ElementHandleArgs) error {
	param := pa.(*image)
	core.SendJSON(p.Res, http.StatusOK, param)
	return nil
}

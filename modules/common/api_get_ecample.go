package common

import (
	"net/http"
	"stbweb/core"
)

type appExample struct{}

func init() {
	core.RegisterFun("example", new(appExample))
}
func (ap *appExample) Get(arge *core.ElementHandleArgs) {
	if arge.APIInterceptionGet("example", nil, appExamplef) {
		return
	}
}
func appExamplef(pa interface{}, content *core.ElementHandleArgs) error {
	core.SendJSON(content.Res, http.StatusOK, core.SendMap{"msg": "resf success"})
	return nil
}

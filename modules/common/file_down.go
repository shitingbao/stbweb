package common

import (
	"net/http"
	"os"
	"path/filepath"
	"stbweb/core"
)

type downFile struct {
	Base string
}

func init() {
	core.RegisterFun("down", new(downFile), true) //??
}

func (d *downFile) Post(p *core.ElementHandleArgs) {
	if p.APIInterceptionPost("down", new(downFile), getFile) {
		return
	}
}

func getFile(param interface{}, p *core.ElementHandleArgs) error {
	pm := param.(*downFile)
	str, err := os.Getwd()
	if err != nil {
		return err
	}
	http.ServeFile(p.Res, p.Req, filepath.Join(str, pm.Base))
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": true})
	return nil
}

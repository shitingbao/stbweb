package common

import (
	"net/http"
	"os"
	"path/filepath"
	"stbweb/core"
)

type downFile struct {
	Base string `json:"base"`
}

func init() {
	core.RegisterFun("down", new(downFile), true)
}

func (d *downFile) Post(p *core.ElementHandleArgs) {
	if p.APIInterceptionPost("down", nil, getFile) {
		return
	}
}

func getFile(param interface{}, p *core.ElementHandleArgs) error {
	pm := param.(*downFile)
	str, err := os.Getwd()
	if err != nil {
		return err
	}
	http.ServeFile(p.Res, p.Req, filepath.Join(str, "assets", pm.Base))
	return nil
}

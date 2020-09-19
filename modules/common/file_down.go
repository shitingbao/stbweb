package common

import (
	"log"
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
	log.Println("getfile")
	pm := param.(*downFile)
	str, err := os.Getwd()
	if err != nil {
		return err
	}
	log.Println("path:", filepath.Join(str, pm.Base))
	http.ServeFile(p.Res, p.Req, filepath.Join(str, pm.Base))
	core.SendJSON(p.Res, http.StatusOK, "success")
	return nil
}

package common

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"stbweb/core"
)

type DownFile struct {
	Base string
}

func init() {
	core.RegisterFun("down", new(DownFile), true) //??
}

func (d *DownFile) Post(p *core.ElementHandleArgs) {
	if p.APIInterceptionPost("down", new(DownFile), getFile) {
		return
	}
}

func getFile(param interface{}, p *core.ElementHandleArgs) error {
	log.Println("getfile")
	pm := param.(*DownFile)
	str, err := os.Getwd()
	if err != nil {
		return err
	}
	log.Println("path:", filepath.Join(str, pm.Base))
	http.ServeFile(p.Res, p.Req, filepath.Join(str, pm.Base))
	core.SendJSON(p.Res, http.StatusOK, "success")
	return nil
}

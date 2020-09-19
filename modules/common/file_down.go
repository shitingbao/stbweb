package common

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"stbweb/core"

	"github.com/sirupsen/logrus"
)

type downFile struct {
	Base string `json:"base"`
}

func init() {
	core.RegisterFun("down", new(downFile), true)
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
		logrus.WithFields(logrus.Fields{"err": err}).Error("down file")
		return err
	}
	log.Println(filepath.Join(str, "assets", pm.Base))
	http.ServeFile(p.Res, p.Req, filepath.Join(str, "assets", pm.Base))
	core.SendJSON(p.Res, http.StatusOK, "success")
	return nil
}

package common

import (
	"net/http"
	"stbweb/core"
	"stbweb/lib/spider"
)

type spiderHand struct{}

func init() {
	core.RegisterFun("spider", new(spiderHand), false)
}

type spiderParam struct {
	URL string
}

func (r *spiderHand) Post(p *core.ElementHandleArgs) {
	if p.APIInterceptionPost("spider", new(spiderParam), spiderLoad) {
		return
	}
}

func spiderLoad(param interface{}, p *core.ElementHandleArgs) error {
	pa := param.(*spiderParam)
	if err := spider.SpiderRun(pa.URL); err != nil {
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"false": err.Error()})
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": true})
	return nil
}

package common

import (
	"stbweb/core"
	"stbweb/lib/rediser"
)

type loginSignOut struct{}

func init() {
	core.RegisterFun("loginout", new(loginSignOut), true)
}

func (l *loginSignOut) Get(p *core.ElementHandleArgs) {
	if p.APIInterceptionGet("loginout", nil, loginOut) {
		return
	}
}

func loginOut(param interface{}, p *core.ElementHandleArgs) error {
	rediser.DelUser(core.Rds, p.Req.Header.Get("token"))
	return nil
}

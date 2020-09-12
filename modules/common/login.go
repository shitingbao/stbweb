package common

import (
	"net/http"
	"stbweb/core"
	"stbweb/lib/rediser"

	"github.com/pborman/uuid"
)

type login struct{}

func init() {
	core.RegisterFun("login", new(login), false)
}

func (l *login) Post(p *core.ElementHandleArgs) {
	if p.APIInterceptionPost("login", new(apiUser), loginAPI) || p.APIInterceptionPost("check", new(tp), checkLogin) {
		return
	}
}

type apiUser struct {
	Name   string
	Pwd    string
	Avatar string
	Email  string
	Phone  string
	Salt   string
	// UpdateTime string
}
type tp struct {
	Token string `json:"token"`
}

func loginAPI(param interface{}, p *core.ElementHandleArgs) error {
	pa := param.(*apiUser)
	if pa.Name == "" || pa.Pwd == "" {
		core.SendJSON(p.Res, http.StatusOK, "必填内容不能为空")
		return nil
	}
	if !core.IsExistUser(pa.Name) {
		core.SendJSON(p.Res, http.StatusOK, "用户不存在")
		return nil
	}
	u := core.GetUser(pa.Name)

	if u.Equal(pa.Pwd) {
		token := uuid.NewUUID().String()
		rediser.RegisterUser(core.Rds, token, pa.Name)
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"token": token, "success": true})
		return nil
	}
	core.SendJSON(p.Res, http.StatusOK, false)
	return nil
}

//检查是否在登录状态
func checkLogin(param interface{}, p *core.ElementHandleArgs) error {
	pa := param.(*tp)
	res := false
	if rediser.CheckLoginUser(core.Rds, pa.Token) {
		res = true
	}
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": res})
	return nil
}

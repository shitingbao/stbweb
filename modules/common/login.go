package common

import (
	"net/http"
	"stbweb/core"
	"stbweb/lib/rediser"

	"github.com/pborman/uuid"
)

type login struct{}

func init() {
	core.RegisterFun("login", new(login))
}
func (l *login) Post(p *core.ElementHandleArgs) {
	if p.APIInterceptionPost("login", new(apiUser), loginAPI) {
		return
	}
}

func (l *login) Get(p *core.ElementHandleArgs) {

}

type apiUser struct {
	Name       string
	Pwd        string
	Avatar     string
	Email      string
	Phone      string
	Salt       string
	UpdateTime string
}

func loginAPI(param interface{}, p *core.ElementHandleArgs) error {
	pa := param.(*apiUser)
	if pa.Name == "" || pa.Pwd == "" {
		core.SendJSON(p.Res, http.StatusOK, "必填内容不能为空")
		return nil
	}
	if !isExistUser(pa.Name) {
		core.SendJSON(p.Res, http.StatusOK, "用户不存在")
		return nil
	}
	u := getUser(pa.Name)

	if u.equal(pa.Pwd) {
		token := uuid.NewUUID().String()
		rediser.RegisterUser(core.Rds, token, pa.Name)
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"token": token, "success": true})
		return nil
	}
	core.SendJSON(p.Res, http.StatusOK, false)
	return nil
}

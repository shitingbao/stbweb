package common

import (
	"net/http"
	"stbweb/core"
	"stbweb/lib/rediser"
)

type login struct{}

type loginParam struct {
	Name string
	Pwd  string
}

func init() {
	core.RegisterFun("login", new(login))
}
func (l *login) Post(p *core.ElementHandleArgs) {
	if p.APIInterceptionPost("login", new(loginParam), loginAPI) {
		return
	}

}

func (l *login) Get(p *core.ElementHandleArgs) {

}

type user struct {
	Name string
	Pwd  string
}

func loginAPI(param interface{}, p *core.ElementHandleArgs) error {
	pa := param.(*loginParam)
	u := user{}

	if err := core.Ddb.QueryRow("select password from user where name=?", pa.Name).Scan(&u.Pwd); err != nil {
		return err
	}
	if pa.Pwd == u.Pwd {
		rediser.RegisterUser(core.Rds, "随机唯一标识码", pa.Name)
		core.SendJSON(p.Res, http.StatusOK, true)
		return nil
	}
	core.SendJSON(p.Res, http.StatusOK, false)
	return nil
}

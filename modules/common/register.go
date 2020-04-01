package common

import (
	"net/http"
	"stbweb/core"
)

type register struct{}

func init() {
	core.RegisterFun("register", new(login))
}

type registerParam struct {
	Name     string
	Email    string
	Phone    string
	Password string
	Avatar   string
}

func (r *register) Post(p *core.ElementHandleArgs) {
	if p.APIInterceptionPost("register", new(registerParam), userRegister) {
		return
	}
}

func userRegister(param interface{}, p *core.ElementHandleArgs) error {
	pa := param.(*registerParam)
	if pa.Name == "" || pa.Password == "" {
		core.SendJSON(p.Res, http.StatusOK, "必填内容不能为空")
		return nil
	}
	if isExistUser(pa.Name) {
		core.SendJSON(p.Res, http.StatusOK, "用户已存在")
		return nil
	}
	sql := "INSERT INTO user(NAME,PASSWORD,avatar,email,phone,salt)VALUES(?,?,?,?,?,?)"
	//两次加密一次解密，双向加单向
	salt := buildIserSalt(pa.Name)
	bPwd := buildPas(pa.Password, salt)
	if _, err := core.Ddb.Exec(sql, pa.Name, string(bPwd), pa.Avatar, pa.Email, pa.Phone, salt); err != nil {
		core.SendJSON(p.Res, http.StatusOK, err.Error())
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, true)
	return nil
}

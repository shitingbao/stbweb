package core

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

var (
	//WebAPIHanderName hander中存放api名称
	WebAPIHanderName = "stbweb-api"
)

//APIInterceptionGet 拦截api请求，统一处理
func (e *ElementHandleArgs) APIInterceptionGet(methodName string, param interface{},
	cb func(pa interface{}, content *ElementHandleArgs) error) bool {
	//名称对应判断
	if methodName != e.apiName() {
		return false
	}
	if err := cb(param, e); err != nil {
		logrus.WithFields(logrus.Fields{
			"elename": e.Element.Name,
			"method":  methodName,
		}).Error(err)
		// SendJSON(e.Res, http.StatusBadRequest, SendMap{"msg": err})
	}
	return true
}

//APIInterceptionPost 拦截api请求，统一处理
//内部需要判断比get多一个body内容不为空和接收类型不为空的判断
func (e *ElementHandleArgs) APIInterceptionPost(methodName string, param interface{},
	cb func(pa interface{}, content *ElementHandleArgs) error) bool {
	//名称对应判断
	if methodName != e.apiName() {
		return false
	}
	if e.Req.ContentLength <= 0 || param == nil {
		logrus.WithFields(logrus.Fields{"param": "nil"}).Error("post")
		return false
	}
	defer e.Req.Body.Close()
	if err := json.NewDecoder(e.Req.Body).Decode(param); err != nil {
		logrus.WithFields(logrus.Fields{"methodName": methodName, "elementName": e.Element.Name}).Error(err)
		return false
	}
	if err := cb(param, e); err != nil {
		logrus.WithFields(logrus.Fields{
			"elename": e.Element.Name,
			"method":  methodName,
		}).Error(err.Error())
	}
	return true
}

func (e *ElementHandleArgs) apiName() string {
	return e.Req.Header.Get(WebAPIHanderName)
}

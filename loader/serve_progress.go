package loader

import (
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"stbweb/core"
	"strings"

	"github.com/Sirupsen/logrus"
)

func init() {
	mime.AddExtensionType(".js", "text/javascript")
	mime.AddExtensionType(".css", "text/css; charset=utf-8")
	// gob.Register(map[string]interface{}{})
}

func httpProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		if core.WebConfig.AllowCORS {
			allowOrigin := core.WebConfig.AllowOrigin
			if len(allowOrigin) == 0 {
				allowOrigin = "*" //待定，跨域允许的指定地址
			}
			w.Header().Set("Access-Control-Allow-Origin", allowOrigin) //设置允许跨域的请求地址
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", fmt.Sprintf(
				"%s,Origin, X-Requested-With, Content-Type, Accept, Connection, User-Agent, Cookie",
				core.WebAPIHanderName)) //这里可以增加对应handle
		}
		logrus.WithFields(logrus.Fields{
			"url":       r.URL.String(),
			"allowCORS": core.WebConfig.AllowCORS,
		}).Warn("options")

		w.WriteHeader(http.StatusOK)
		return
	}
	//？？这里需要定向前端地址，待定
	if r.URL.String() == "/" {
		// http.ServeFile(w, r, filepath.Join("./dist", "index.html"))
		str, err := os.Getwd()
		if err != nil {
			logrus.WithFields(logrus.Fields{"path": err.Error()}).Warn("getwd")
			return
		}
		http.ServeFile(w, r, filepath.Join(str, "dist", "index.html"))
		return
	}
	paths, err := parsePaths(r.URL)
	//这里的path反馈工作元素内容待定
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(nil)
		return
	}
	suf := strings.Split(paths[len(paths)-1], ".")
	if suf[len(suf)-1] == "css" {
		str, err := os.Getwd()
		if err != nil {
			logrus.WithFields(logrus.Fields{"path": err.Error()}).Warn("getwd")
			return
		}
		http.ServeFile(w, r, filepath.Join(str, "dist", "css", paths[len(paths)-1]))
		return
	}
	if suf[len(suf)-1] == "js" {
		str, err := os.Getwd()
		if err != nil {
			logrus.WithFields(logrus.Fields{"path": err.Error()}).Warn("getwd")
			return
		}
		http.ServeFile(w, r, filepath.Join(str, "dist", "js", paths[len(paths)-1]))
		return
	}
	core.ElementHandle(w, r, paths[0]) //待定，工作元素的名称获取是否来源于路由
}
func fileHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "./dist/index.html", http.StatusFound)
}

func parsePaths(u *url.URL) ([]string, error) {
	paths := []string{}
	pstr := u.EscapedPath()

	for _, str := range strings.Split(pstr, "/")[1:] {
		s, err := url.PathUnescape(str)
		if err != nil {
			return nil, err
		}
		paths = append(paths, s)
	}
	return paths, nil
}

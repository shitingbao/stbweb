package loader

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"stbweb/core"
	"strings"

	"github.com/Sirupsen/logrus"
)

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
	// http.Handle("/dist", http.StripPrefix("/dist", http.FileServer(http.Dir("dist"))))
	//？？这里需要定向前端地址，待定
	if r.URL.String() == "/" {
		// http.Redirect(w, r, "dist/index.html", http.StatusFound)
		// http.ServeFile(w, r, filepath.Join("dist", "index.html")) //配置自己的前端入口，找不到会404
		// core.SendJSON(w, http.StatusOK, core.SendMap{"url": "nothing"})

		f, err := os.Open("/opt/dist/index.html")
		if err != nil {
			log.Println(err)
			return
		}
		defer f.Close()

		data, err := ioutil.ReadAll(f)
		if err != nil {
			log.Println(err)
			return

		}
		w.Write(data)
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
	core.ElementHandle(w, r, paths[0]) //待定，工作元素的名称获取是否来源于路由
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

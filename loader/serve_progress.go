package loader

import (
	"net/http"
	"net/url"
	"stbweb/core"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func httpProcess(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() == "/" {
		core.SendJSON(w, http.StatusOK, core.SendMap{"url": "nothing"})
		return
	}
	paths, err := parsePaths(r.URL)
	//这里的path反馈工作元素内容待定
	if err != nil {
		core.LOG.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(nil)
		return
	}
	core.LOG.WithFields(log.Fields{
		"paths": paths,
	}).Info("paths")
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

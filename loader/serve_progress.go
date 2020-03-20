package loader

import (
	"net/http"
	"net/url"
	"stbweb/core"
	"strings"
)

func httpProcess(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() == "/" {
		//待定，可以反馈静态资源或者文档地址
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
	core.ElementHandle(w, r, paths[0])
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

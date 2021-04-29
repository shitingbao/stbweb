package spider

import (
	"io"
	"net/http"
	"os"
	"path"
	"stbweb/core"
	"strings"
)

var (
	img         = "https://www.pximg.com/wp-content/uploads/2020/04/cdd72379d64dd08-scaled.jpg"
	url         = "https://www.pximg.com/meinv/39211.html"
	baseFileDir = core.DefaultFilePath
)

func urlQueryHandle(imageURL string) string {
	if imageURL == "" {
		return ""
	}
	s := strings.Split(imageURL, "&")
	return s[0]
}

// 路径中可能带？的参数，需要处理一下
func createImage(imageURL string) error {
	resp, err := http.Get(imageURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return err
	}
	bt, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	nakedURL := urlQueryHandle(imageURL)
	p := setGetFileDir(nakedURL)
	fileName := path.Join(p, path.Base(nakedURL))
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	if _, err := f.Write(bt); err != nil {
		return err
	}
	return nil
}

// 每一个页面一个文件夹，取自倒数第二个目录名称,并建立目录
func setGetFileDir(imageURL string) string {
	p := path.Join(baseFileDir, path.Base(path.Dir(imageURL)))
	os.MkdirAll(p, os.ModePerm)
	return p
}

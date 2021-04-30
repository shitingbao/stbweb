package spider

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"stbweb/core"
	"strings"
)

var (
	baseFileDir = core.DefaultFilePath
)

// 去除名称中一些特殊符号
func fileNameHandle(n string) string {
	n = strings.ReplaceAll(n, " ", "")
	n = strings.ReplaceAll(n, "-", "_")
	n = strings.ReplaceAll(n, ".", "_")
	n = strings.ReplaceAll(n, "!", "_")
	return n
}

// 路径中可能带？的参数，需要处理一下
func createImage(imageURL, fBaseName string) error {

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
	fileName := getFileName(imageURL, fBaseName)

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
func bashPath(imageURL string) string {
	p := path.Join(baseFileDir, path.Base(path.Dir(imageURL)))
	os.MkdirAll(p, os.ModePerm)
	return p
}

func getFileName(imageURL, fBaseName string) string {
	u, err := url.Parse(imageURL)
	if err != nil {
		return ""
	}
	bPath := u.Query().Get("src")
	p := bashPath(bPath)
	return path.Join(p, fBaseName+path.Base(bPath))
}

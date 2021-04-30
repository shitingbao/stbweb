package spider

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"golang.org/x/net/html"
)

type imgNode struct {
	Attr     []html.Attribute
	Url      *url.URL
	Resp     *http.Response
	FileName string
	Src      string
}

func NewImgNode(attr []html.Attribute, resp *http.Response) *imgNode {
	return &imgNode{
		Attr: attr,
		Resp: resp,
	}
}

func (i *imgNode) Handle() error {
	for _, a := range i.Attr {
		switch {
		case a.Key == "src": //进入这两个说明是图片标签
			link, err := i.Resp.Request.URL.Parse(a.Val)
			if err != nil {
				return err
			}
			i.Src = link.String()
		case a.Key == "alt":
			i.FileName = fileNameHandle(a.Val)
		}
	}
	return i.createImage()
}

// 路径中可能带？的参数，需要处理一下
func (i *imgNode) createImage() error {
	resp, err := http.Get(i.Src)
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
	fileName := i.getFileName()
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
func (i *imgNode) bashPathMkDir() string {
	p := path.Join(baseFileDir, path.Base(path.Dir(i.Src)))
	os.MkdirAll(p, os.ModePerm)
	return p
}

func (i *imgNode) getFileName() string {
	u, err := url.Parse(i.Src)
	if err != nil {
		return ""
	}
	bPath := u.Query().Get("src")
	if bPath == "" {
		bPath = i.Src
	}
	p := i.bashPathMkDir()
	suffixPath := path.Base(bPath)
	lasPath := strings.Split(suffixPath, "?") //可能还有参数
	return path.Join(p, i.FileName+lasPath[0])
}

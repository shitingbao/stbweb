package spider

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"stbweb/core"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

var (
	baseFileDir = core.DefaultFilePath
)

// imgNode img标签对象
// Attr标签属性切片
// Url
// Resp该标签页面的请求对象
// FileName要保存的文件名称
// Src img对应的src属性内容地址，也就是图像地址
type imgNode struct {
	FileName string
	Src      string
	Node     *html.Node
	Resp     *http.Response
	c        chan *imgNode
}

func NewImgNode(resp *http.Response, n *html.Node) htmlNode {
	return &imgNode{
		Node: n,
		Resp: resp,
	}
}

// 进入异步准备
func (i *imgNode) Handle(ch chan *imgNode) error {
	if i.Node.Data != imgNodeSign {
		return nil
	}
	for _, a := range i.Node.Attr {
		switch {
		case a.Key == srcNodeSign: // 进入这两个说明是图片标签
			link, err := i.Resp.Request.URL.Parse(a.Val)
			if err != nil {
				return err
			}
			i.Src = link.String()
		case a.Key == altNodeSign:
			i.fileNameHandle(a.Val)
		}
	}
	ch <- i
	// return i.createImage()
	return nil
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

// 去除名称中一些特殊符号
func (i *imgNode) fileNameHandle(n string) {
	n = strings.ReplaceAll(n, " ", "")
	n = strings.ReplaceAll(n, "-", "_")
	n = strings.ReplaceAll(n, ".", "_")
	n = strings.ReplaceAll(n, "!", "_")
	i.FileName = n
}

func imageConsumer(ctx context.Context, ch <-chan *imgNode) {
	for {
		select {
		case <-ctx.Done():
			logrus.Info("out create")
			return
		case i, ok := <-ch:
			if !ok {
				return
			}
			if i != nil {
				i.createImage()
			}
		}
	}
}

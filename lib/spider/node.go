package spider

// 注意，不能直接对每个页面进行异步操作
// 假设每个页面加一个异步携程，每个页面下都可能有不定量的新页面，这时候如果携程池固定协程数量，那后面的就在等待
// 同时，前面的页面中开了一半（比如有20个新页面，但是只有10个协程），这时候就卡死了
import (
	"context"
	"net/http"
	"runtime"

	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

var (
	nodeSign     = "node"
	aNodeSign    = "a"
	imgNodeSign  = "img"
	hrefNodeSign = "href"
	srcNodeSign  = "src"
	altNodeSign  = "alt"
)

// forEachNode 便利一个node中所有的节点（就是所有标签）
// 一个页面中的所有标签便利
func forEachNode(resp *http.Response, n *html.Node, ch chan *imgNode) error {
	if err := nodefunc(resp, n, ch); err != nil {
		logrus.WithFields(logrus.Fields{"node": n.Data, "err": err.Error()}).Error("nodefunc")
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if err := forEachNode(resp, c, ch); err != nil {
			return err
		}
	}
	return nil
}

// 一个标签内的处理
func nodefunc(resp *http.Response, n *html.Node, ch chan *imgNode) error { //页面中一个标签单次节点处理
	if n.Type != html.ElementNode {
		return nil
	}
	var h htmlNode
	switch n.Data {
	case aNodeSign:
		h = NewANode(resp, n)
	case imgNodeSign:
		h = NewImgNode(resp, n)
	default:
		return nil
	}
	return h.Handle(ch)
}

// GetURLInfomationAdress is a func get URL infomation
func SpiderLoad(URL string, ch chan *imgNode) error {
	resp, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return err
	}
	n, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}
	forEachNode(resp, n, ch)
	return nil
}

// nodeSign 作为一个的单次的标识数，以免多次提交后跳过
func SpiderRun(URL string) error {
	nodeSign = uuid.New()
	ch := make(chan *imgNode)
	ctx, canlf := context.WithCancel(context.Background())
	defer canlf()

	for i := 0; i < runtime.NumCPU()+1; i++ {
		go imageConsumer(ctx, ch)
	}
	return SpiderLoad(URL, ch)
}

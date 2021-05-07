package spider

// 放弃异步处理每个页面，无法判定遍历完成所有的网站节点

import (
	"net/http"
	"stbweb/core"
	"sync"

	"golang.org/x/net/html"
)

// 用于去重检查
var (
	pageMap  = sync.Map{}
	newsChan = "html_news_chan"
)

type img struct {
	Src string
	Alt string
}

// 用于异步的节点对象
// 过程：
// pageMap中key检查
// lock，double check
// pageMap加入新url
// 新url同时放入redis 队列（集合）中等待下次处理
type nodeSync struct {
	mu *sync.Mutex
}

func (node *nodeSync) forEachNodeSync(resp *http.Response, n *html.Node) error {
	imgUrls := []img{}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		switch n.Data {
		case aNodeSign:
			for _, at := range n.Attr {
				if at.Key != hrefNodeSign {
					continue
				}
				link, err := resp.Request.URL.Parse(at.Val)
				if err != nil {
					return err
				}
				l := link.String()
				if _, ok := pageMap.Load(l); ok {
					continue
				}
				// doubke check
				// into redis queue
				// Uniqueness check
				node.mu.Lock()
				if _, ok := pageMap.Load(l); !ok {
					pageMap.Store(l, "1")
					core.Rds.SAdd(newsChan, l)
				}
				node.mu.Unlock()
			}
		case imgNodeSign:
			i := img{}
			for _, attr := range n.Attr {
				switch {
				case attr.Key == srcNodeSign: // 进入这两个说明是图片标签
					link, err := resp.Request.URL.Parse(attr.Val)
					if err != nil {
						return err
					}
					i.Src = link.String()
				case attr.Key == altNodeSign:
					i.Alt = attr.Val
				}
			}
			imgUrls = append(imgUrls, i)
		}
	}
	return nil
}

func getNews() (string, error) {
	return core.Rds.SPop(newsChan).Result()
}

// func spiderLoadSync() error {
// 	for {
// 		l, err := getNews()
// 		if err != nil {
// 			return err
// 		}
// 		resp, err := http.Get(l)
// 		if err != nil {
// 			return err
// 		}
// 		defer resp.Body.Close()
// 		if resp.StatusCode != http.StatusOK {
// 			return err
// 		}
// 		n, err := html.Parse(resp.Body)
// 		if err != nil {
// 			return err
// 		}
// 		forEachNode(resp, n)
// 	}
// }

package spider

import (
	"log"
	"net/http"
	"stbweb/core"

	"golang.org/x/net/html"
)

var (
	nodeSign = "node"
)

// forEachNode 便利一个node中所有的节点（就是所有标签）
func forEachNode(resp *http.Response, n *html.Node) error {
	// core.WorkPool.Submit(func() {
	// 	if err := nodefunc(resp, n); err != nil {
	// 		log.Println(err)
	// 	}
	// })
	if err := nodefunc(resp, n); err != nil {
		log.Println(err)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if err := forEachNode(resp, c); err != nil {
			return err
		}
	}
	return nil
}

func nodefunc(resp *http.Response, n *html.Node) error { //页面中一个标签单次节点处理
	if n.Type != html.ElementNode || (n.Data != "a" && n.Data != "img") {
		return nil
	}

	for _, a := range n.Attr {
		link, err := resp.Request.URL.Parse(a.Val)
		if err != nil {
			return err
		}
		l := link.String()
		switch a.Key {
		case "href":
			if core.Rds.HGet(nodeSign, l).Val() != "" {
				continue
			}
			core.Rds.HSet(nodeSign, l, nodeSign)
			SpiderLoad(a.Val)
		case "src":
			if err := createImage(l); err != nil {
				return err
			}
		}
	}
	return nil
}

//GetURLInfomationAdress is a func get URL infomation
func SpiderLoad(URL string) error {

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
	forEachNode(resp, n)
	return nil
}

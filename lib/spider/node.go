package spider

import (
	"log"
	"net/http"
	"stbweb/core"
	"sync"
	"time"

	"golang.org/x/net/html"
)

type HtmlNodeHandle struct {
	Url string
	Mu  *sync.Mutex
}

func (h *HtmlNodeHandle) getNodeKey(url string) bool {
	if core.Rds.HGet(url, url).Val() == "" {
		return true
	}
	return false
}

func (h *HtmlNodeHandle) setNodeKey(url string) bool {
	if !h.getNodeKey(url) {
		return false
	}
	h.Mu.Lock()
	if !h.getNodeKey(url) { //double check
		return false
	}
	core.Rds.HSet(url, url, url)
	return true
}

// ForEachNode 便利一个node中所有的节点（就是所有标签）
func ForEachNode(resp *http.Response, n *html.Node, f func(resp *http.Response, n *html.Node)) {
	if f != nil {
		f(resp, n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ForEachNode(resp, c, f)
	}
}

func ForNodeElementfunc(resp *http.Response, n *html.Node) { //单次节点处理
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key != "href" {
				continue
			}
			link, err := resp.Request.URL.Parse(a.Val)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("link")

			core.Rds.SetNX(link.String(), "", time.Second)
			// if CheckURL(link.String(), links) {// 检查重复。setnx
			// 	links = append(links, link.String()) //这条语句可以改成并行获取URL地址内容
			// }
		}
	}
}

//GetURLInfomationAdress is a func get URL infomation
func GetURLInfomationAdress(URL string) []string {
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatal("Can't connect:", URL)
	}
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var links []string

	ForEachNode(resp, doc, ForNodeElementfunc)
	return links
}

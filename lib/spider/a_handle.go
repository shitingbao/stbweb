package spider

import (
	"golang.org/x/net/html"
)

type aNode struct {
	Attr []html.Attribute
	Src  string
}

func NewANode(attr []html.Attribute) *aNode {
	return &aNode{
		Attr: attr,
	}
}

func (a *aNode) Hand() {}

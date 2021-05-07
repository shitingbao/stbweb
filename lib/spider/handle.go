package spider

type htmlNode interface {
	Handle(ch chan *imgNode) error
}

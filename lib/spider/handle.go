package spider

type htmlNode interface {
	Handle() error
}

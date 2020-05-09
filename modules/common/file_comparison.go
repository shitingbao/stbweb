package common

import (
	"stbweb/core"
)

type fileComparison struct{}

func init() {
	core.RegisterFun("filecomparison", new(fileComparison), false)
}

func (f *fileComparison) Get(p *core.ElementHandleArgs) {
	if p.APIInterceptionGet("comparison", nil, comparison) {
		return
	}
}

func comparison(param interface{}, p *core.ElementHandleArgs) error {
	return nil
}

func (f *fileComparison) Post(p *core.ElementHandleArgs) {}

package common

import (
	"stbweb/core"
	"stbweb/lib/comparison"
)

type fileComparison struct{}

func init() {
	core.RegisterFun("filecomparison", new(fileComparison), false)
}

func (f *fileComparison) Get(p *core.ElementHandleArgs) {
	if p.APIInterceptionGet("comparison", nil, comparisonFile) {
		return
	}
}

func comparisonFile(param interface{}, p *core.ElementHandleArgs) error {
	name := p.Req.URL.Query().Get("filename")
	sep := p.Req.URL.Query().Get("sep")
	// comparison.GetTitleLineGroup(name, sep)
	comparison.GetLineGroup(name, sep)
	// comparison.ExcelTitleLineGroup(name)
	// comparison.ExcelLineGroup(name)
	return nil
}

func (f *fileComparison) Post(p *core.ElementHandleArgs) {}

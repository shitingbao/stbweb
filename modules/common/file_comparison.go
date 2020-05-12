package common

import (
	"net/http"
	"stbweb/core"
	"stbweb/lib/comparison"
)

type fileComparison struct{}

func init() {
	core.RegisterFun("filecomparison", new(fileComparison), false)
}

type comparisonParam struct {
	CompareFile comparison.ParisonFileObject
	AimFile     comparison.ParisonFileObject
}

func (f *fileComparison) Post(p *core.ElementHandleArgs) {
	if p.APIInterceptionPost("comparison", new(comparisonParam), comparisonFileCommon) {
		return
	}
}

func comparisonFileCommon(param interface{}, p *core.ElementHandleArgs) error {
	pa := param.(*comparisonParam)
	res, err := comparison.FileComparison(pa.CompareFile, pa.AimFile)
	if err != nil {
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, res)
	return nil
}

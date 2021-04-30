package spider

import (
	"stbweb/core"
	"strings"
)

var (
	baseFileDir = core.DefaultFilePath
)

// 去除名称中一些特殊符号
func fileNameHandle(n string) string {
	n = strings.ReplaceAll(n, " ", "")
	n = strings.ReplaceAll(n, "-", "_")
	n = strings.ReplaceAll(n, ".", "_")
	n = strings.ReplaceAll(n, "!", "_")
	return n
}

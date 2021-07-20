package math

// 绝对值
func CalcAbs(a int) (ret int) {
	ret = (a ^ a>>31) - a>>31
	return
}

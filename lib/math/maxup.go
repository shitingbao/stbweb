package math

import "sort"

// list 中求最大递增的子序列的长度，只是长度
func lengthOfLIS(nums []int) int {
	dlist := []int{}
	for _, v := range nums {
		idx := sort.SearchInts(dlist, v)
		if idx < len(dlist) {
			dlist[idx] = v
		} else {
			dlist = append(dlist, v)
		}
	}
	return len(dlist)
}

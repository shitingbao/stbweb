package math

// 二分法，在 nums1 中找到最接近 flag 的值，反馈索引
func two(nums1 []int, flag int) int {
	left := 0
	right := len(nums1)
	idx := len(nums1) / 2
	for {
		if idx == 0 {
			return idx
		}
		if idx == len(nums1)-1 {
			if abs(nums1[idx], flag) < abs(nums1[idx-1], flag) {
				return idx
			} else {
				return idx - 1
			}

		}
		if flag == nums1[idx] {
			break
		}
		if flag > nums1[idx-1] && flag < nums1[idx] {
			if abs(nums1[idx], flag) > abs(nums1[idx-1], flag) {
				idx--
			}
			break
		}
		if flag > nums1[idx] {
			left = idx
			idx = left + (right-left)/2
		} else {
			right = idx
			idx = left + (right-left)/2
		}
	}
	return idx
}

func abs(a, b int) int {
	if a > b {
		return a - b
	} else {
		return b - a
	}
}

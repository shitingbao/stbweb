package main

//反馈列表中两个数值之间，大的数进行保存
func getMaxArea(heights []int) int {
	resArea := 0
	skey := 0
	for i, v := range heights {
		if v == 0 {
			area := getExcelArea(heights[skey:i])
			if area > resArea {
				resArea = area
			}
			skey = i + 1
			continue
		}
		area := getExcelArea(heights[skey:])
		if area > resArea {
			resArea = area
		}
	}
	return resArea
}

//返回不为0列表的最大面积
func getExcelArea(sli []int) int {
	if len(sli) == 1 {
		return sli[0]
	}
	res := 0
	for i, v := range sli {
		if v > res {
			res = v
		}
		lsum := 0
		rsum := 0
		//向左找到最低值
		for lkey := i; lkey >= 0; lkey-- {
			if sli[lkey] < v {
				break
			}
			lsum++
		}
		//向右找到最低值
		for rkey := i; rkey < len(sli); rkey++ {
			if sli[rkey] < v {
				break
			}
			rsum++
		}
		sum := lsum + rsum
		if sum != 0 {
			sum--
		}
		if sum*v > res {
			res = sum * v
		}
	}
	return res
}

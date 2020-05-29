package main

//ali面试题，求水滴
//传入一组int数组，反馈中间为0被不为0包含的数量,int值为高度，最大值就是最高高度，没层一次减一
//如：[1,0,1]，输出1，[1,0,0,1,0,1]，输出3，[1, 0, 2, 0, 0, 1, 3, 0, 1]输出7
type slice []int

//反馈水滴数量
func getResCount(sli []int) int {
	count := 0
	lines := splitSli(sli)
	for _, v := range lines {
		count += getCount(v)
	}
	return count
}

func getSliMax(sli []int) int {
	count := 0
	for _, v := range sli {
		if v > count {
			count = v
		}
	}
	return count
}

func splitSli(sli []int) map[int]slice {
	max := getSliMax(sli)
	res := make(map[int]slice)
	for i := 0; i < max; i++ {
		lx := []int{}
		for _, v := range sli {
			if v > 0 && v > i {
				lx = append(lx, v-i)
				continue
			}
			lx = append(lx, 0)
		}
		res[i+1] = lx
	}
	return res
}

func getCount(sli []int) int {
	count := 0
	key := -1
	resCount := 0
	for i := range sli {
		if i != 0 && sli[i] != 0 && sli[i-1] == 0 && key != -1 {
			count = i - key - 1
			// log.Println(key, "--", i, "--", count)
			resCount += count
			count = 0
			key = -1
		}
		if i != len(sli)-1 && sli[i] != 0 && sli[i+1] == 0 {
			key = i
		}
	}
	return resCount
}

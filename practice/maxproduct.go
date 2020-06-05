package main

func maxProduct(words []string) int {

	max := 0
	key := []map[string]int{}
	for idv, v := range words {
		tp := make(map[string]int)
		// isNum := true
		for _, val := range v {
			tp[string(val)]++
			if idv == 0 {
				continue
			}
			for i := 0; i < idv; i++ {
				if key[i][string(val)] != 0 {
					// isNum = false
					break
				}
			}

			// if isNum {
			// 	if len(v)*len(words[i]) > max {
			// 		max = len(v) * len(words[i])
			// 	}
			// }
		}

		key = append(key, tp)
	}
	return max
}

package main

//给定一个字符串，请你找出其中不含有重复字符的 最长子串 的长度。
func lengthOfLongestSubstring(s string) int {
	switch len(s) {
	case 0:
		return 0
	case 1:
		return 1
	}
	max := 1
	keyStr := s[:1]
	for i := 1; i < len(s); i++ {
		keyCount := searchSame(keyStr, string(s[i]))
		if keyCount == -1 {
			keyStr += string(s[i])
			if len(keyStr) > max {
				max = len(keyStr)
			}
			continue
		}
		keyStr = keyStr[keyCount+1:] + string(s[i])
	}
	return max
}

func searchSame(s, sep string) int {
	for i, v := range s[:] {
		if string(v) == sep {
			return i
		}
	}
	return -1
}

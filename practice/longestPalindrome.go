package main

import "log"

//寻找最长的回文子串，中心扩展的方法，就是以一个索引点为中心，向两边扩展比较，注意奇数和偶数
func longestPalindrome(s string) string {
	res := ""
	switch len(s) {
	case 0, 1:
		return s
	case 2:
		if s[0] == s[1] {
			return s
		}
		return s[:1]
	default:
		for i := range s {
			if i == 0 {
				continue
			}
			sp := objecgLeng(i, s)
			if len(sp) > len(res) {
				res = sp
			}
		}
	}
	return res
}

func objecgLeng(logo int, str string) string {
	index := 0
	res := ""
	for i := 0; i <= logo; i++ {
		if logo-i < 0 || logo+i > len(str)-1 {
			break
		}
		if str[logo-i] == str[logo+i] {
			index = i
			res = str[logo-i : logo+i+1]
		} else {
			break
		}
	}
	for i := 0; i <= logo; i++ {
		if logo-i-1 < 0 || logo+i > len(str)-1 {
			break
		}
		if str[logo-i-1] == str[logo+i] {
			if i >= index {
				index = i
				res = str[logo-i-1 : logo+i+1]
			}
		} else {
			break
		}
	}
	return res
}

func longestPalindromeLoad() {
	str := longestPalindrome("cbbd")
	log.Println(str)
}

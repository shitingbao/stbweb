package main

import "log"

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
	res := ""
	isB := true
	for i := 0; i <= logo; i++ {
		if logo-i < 0 || logo+i > len(str)-1 {
			break
		}
		if str[logo-i] != str[logo+i] {
			isB = false
			res = str[logo-i+1 : logo+i]
			break
		}
	}
	for i := 0; i <= logo; i++ {
		if logo-i-1 < 0 || logo+i > len(str)-1 {
			break
		}
		if str[logo-i-1] != str[logo+i] {
			isB = false
			sp := str[logo-i : logo+i]
			if len(sp) > len(res) {
				res = sp
			}
			break
		}
	}
	if isB {
		return str
	}
	return res
}

func longestPalindromeLoad() {
	str := longestPalindrome("ccc")
	log.Println(str)
}

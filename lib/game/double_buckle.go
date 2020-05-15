package game

//DoubleBuckle 双
type DoubleBuckle struct{}

//BrandComparison 比较大小
func (db *DoubleBuckle) BrandComparison(d, aim DeckOfCards) bool {
	return true
}

//GetBrandType 类型反馈
func (db *DoubleBuckle) GetBrandType(d DeckOfCards) string {
	switch len(d.Bd) {
	case 1:
		return Single
	case 2:
		if d.Bd[0] == d.Bd[1] {
			return Double
		}
		return CodeErr
	case 3:
		if d.Bd[0] == d.Bd[1] && d.Bd[1] == d.Bd[2] {
			return Three
		}
		return CodeErr
	default:
	}
	return ""
}

//JudgeStraight 判断是顺子
func JudgeStraight(str []int) bool {
	if len(str) < 5 {
		return false
	}
	for i := 0; i < len(str)-1; i++ {
		if str[i+1]-str[i] != 1 {
			return false
		}
	}
	return true
}

//JudgeEvenPair 连对判断,只要判断到倒数第二和倒数第一相等，倒数第二和倒数第三差1即可
func JudgeEvenPair(str []int) bool {
	if len(str) < 6 && len(str)%2 != 0 {
		return false
	}
	for i := 0; i < len(str)-1; i++ {
		switch {
		case i%2 == 0:
			if str[i] != str[i+1] {
				return false
			}
		default:
			if str[i+1]-str[i] != 1 {
				return false
			}
		}
	}
	return true
}

//JudgeTriplePair 三连对
func JudgeTriplePair(str []int) bool {
	if len(str) >= 9 && len(str)%3 != 0 {
		return false
	}
	for i := 0; i < len(str)-2; i++ {
		switch {
		case i%3 == 0:
			if !(str[i] == str[i+1] && str[i] == str[i+2]) {
				return false
			}
		case i%3 == 1:
			if str[i+2]-str[i] != 1 {
				return false
			}
		default:
		}
	}
	return true
}

//Judgebomb 炸弹
func Judgebomb(str []int) bool {
	if len(str) < 4 {
		return false
	}
	for i := 0; i < len(str)-1; i++ {
		if str[i] != str[i+1] {
			return false
		}
	}
	return true
}

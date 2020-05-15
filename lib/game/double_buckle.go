package game

//DoubleBuckle 双
type DoubleBuckle struct{}

//BrandComparison 比较大小
func (db *DoubleBuckle) BrandComparison(d, aim DeckOfCards) bool {
	if db.GetBrandType(d) != db.GetBrandType(aim) {
		return false
	}
	//炸要另外判断
	return true
}

//GetBrandType 类型反馈
//必须先调用SortOrderAsc排序
func (db *DoubleBuckle) GetBrandType(d DeckOfCards) string {
	d.SortOrderAsc() //必须先进行排序
	switch {
	case len(d.Bd) == 1:
		return Single
	case len(d.Bd) == 2:
		if d.Bd[0] == d.Bd[1] {
			return Double
		}
		return CodeErr
	case len(d.Bd) == 3:
		if d.Bd[0] == d.Bd[1] && d.Bd[1] == d.Bd[2] {
			return Three
		}
		return CodeErr
	case JudgeStraight(d):
		return Straight
	case JudgeEvenPair(d):
		return EvenPair
	case JudgeTriplePair(d):
		return ThreeEvenPair
	case Judgebomb(d):
		return Bomb
	default:
		return CodeErr
	}
}

//JudgeStraight 判断是顺子
func JudgeStraight(str DeckOfCards) bool {
	if len(str.Bd) < 5 {
		return false
	}
	for i := 0; i < len(str.Bd)-1; i++ {
		if str.Bd[i+1].Code-str.Bd[i].Code != 1 {
			return false
		}
	}
	return true
}

//JudgeEvenPair 连对判断,只要判断到倒数第二和倒数第一相等，倒数第二和倒数第三差1即可
func JudgeEvenPair(str DeckOfCards) bool {
	if len(str.Bd) < 6 && len(str.Bd)%2 != 0 {
		return false
	}
	for i := 0; i < len(str.Bd)-1; i++ {
		switch {
		case i%2 == 0:
			if str.Bd[i].Code != str.Bd[i+1].Code {
				return false
			}
		default:
			if str.Bd[i+1].Code-str.Bd[i].Code != 1 {
				return false
			}
		}
	}
	return true
}

//JudgeTriplePair 三连对
func JudgeTriplePair(str DeckOfCards) bool {
	if len(str.Bd) >= 9 && len(str.Bd)%3 != 0 {
		return false
	}
	for i := 0; i < len(str.Bd)-2; i++ {
		switch {
		case i%3 == 0:
			if !(str.Bd[i].Code == str.Bd[i+1].Code && str.Bd[i].Code == str.Bd[i+2].Code) {
				return false
			}
		case i%3 == 1:
			if str.Bd[i+2].Code-str.Bd[i].Code != 1 {
				return false
			}
		default:
		}
	}
	return true
}

//Judgebomb 炸弹
func Judgebomb(str DeckOfCards) bool {
	if len(str.Bd) < 4 {
		return false
	}
	for i := 0; i < len(str.Bd)-1; i++ {
		if str.Bd[i].Code != str.Bd[i+1].Code {
			return false
		}
	}
	return true
}

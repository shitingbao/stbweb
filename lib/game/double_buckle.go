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

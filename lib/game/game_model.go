//Package game 出的阶段,用一个数组，保存未出的其他玩家，在都不要或者其他人出后，更新
package game

//类型规则
const (
	Single   = "Single"
	Double   = "Double"
	Three    = "Three"
	Straight = "Straight" //顺
	EvenPair = "EvenPair" //对
	Bomb     = "Bomb"
	CodeErr  = "Err" //错误类型
)

//Brand 基本
type Brand struct {
	Code int
	Suit string
}

//CodeGameType 类型
type CodeGameType interface {
	BrandComparison(d, aim DeckOfCards) bool
	GetBrandType(d DeckOfCards) string
}

//DeckOfCards 堆
type DeckOfCards struct {
	Bd []Brand
}

//同一类的比较
func stackComparison(cardType string) {
	switch cardType {
	case Single:
	case Double:
	case Three:
	case Straight:
	case EvenPair:
	case Bomb:
	}
}

//调用下列比较函数前需要先校验类型，有类型之后主要就是大小比较，待定

//AllComparison 长度相同大小比较，特殊情况不同，比如炸
func AllComparison(brd, aim DeckOfCards) bool {
	brdNum := 0
	for _, v := range brd.Bd {
		brdNum += v.Code
	}
	aimNum := 0
	for _, v := range brd.Bd {
		aimNum += v.Code
	}
	if brdNum > aimNum {
		return true
	}
	return false
}

//SortOrderAsc 升序排序
func (dc *DeckOfCards) SortOrderAsc() {
	pt := Brand{}
	for i := 0; i < len(dc.Bd); i++ {
		temp := -1
		tempCode := dc.Bd[i].Code
		for j := len(dc.Bd) - 1; j > i; j-- {
			if dc.Bd[j].Code < tempCode {
				temp = j
				tempCode = dc.Bd[j].Code
			}
		}
		if temp != -1 {
			pt = dc.Bd[temp]
			dc.Bd[temp] = dc.Bd[i]
			dc.Bd[i] = pt
		}
	}
}

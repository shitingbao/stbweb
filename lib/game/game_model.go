//Package game 出的阶段,用一个数组，保存未出的其他玩家，在都不要或者其他人出后，更新
package game

//类型规则
const (
	Single        = "单牌"
	Double        = "对牌"
	Three         = "三牌"
	Straight      = "顺子"
	EvenPair      = "连对"
	ThreeEvenPair = "三连对"
	Bomb          = "炸"
	CodeErr       = "Err" //错误类型

	PlumBlossom = "梅花"
	Square      = "方块"
	Spades      = "黑桃"
	RedPeach    = "红桃"
)

var (
	gUser = GameUser{}
)

type GameUser struct {
	Users []string
}

//Brand 基本
type Brand struct {
	Code int
	Suit string
}

//CodeGameType 类型
type CodeGameType interface {
	LicensingCode() DeckOfCards
	BrandComparison(d, aim DeckOfCards) bool
	GetBrandType(d DeckOfCards) string
}

//DeckOfCards 堆
type DeckOfCards struct {
	Bd []Brand
}

//AllComparison 长度相同大小比较，特殊情况不同
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

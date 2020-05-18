package game

import (
	"encoding/json"
	"errors"
	"log"
	"stbweb/lib/ws"
)

//DoubleBuckle 双
type DoubleBuckle struct{}

//AllUserBrands 保存四个用户对应数据
type AllUserBrands struct {
	UserBrands map[string]DeckOfCards
}

var (
	dbuck         DoubleBuckle
	lastBrands    DeckOfCards   //保存上一次的数据
	allUserBrands AllUserBrands //保存四个人的所有数据
)

//LicensingCode 给定所有牌,10以后，11，12，13，14，15，16增量代表J，Q，K，A，2,小王，大王
func (dk *DoubleBuckle) LicensingCode() DeckOfCards {
	dModel := DeckOfCards{}
	suit := ""
	for m := 0; m < 2; m++ {
		switch m {
		case 0:
			suit = PlumBlossom
		case 1:
			suit = Square
		case 2:
			suit = Spades
		case 3:
			suit = RedPeach
		}
		for i := 3; i < 16; i++ {
			bm := Brand{
				Code: i,
				Suit: suit,
			}
			dModel.Bd = append(dModel.Bd, bm)
		}
		bt := Brand{
			Code: 16,
			Suit: "",
		}
		dModel.Bd = append(dModel.Bd, bt)
		bt.Code = 17
		dModel.Bd = append(dModel.Bd, bt)
	}
	return dModel
}

//BrandComparison 比较大小,其中aim为后出的逻辑牌队列
//其中第一个为bomb的情况时，第二个一定为bomb才能比较
func (dk *DoubleBuckle) BrandComparison(d, aim DeckOfCards) bool {
	if (dk.GetBrandType(d) != dk.GetBrandType(aim)) && dk.GetBrandType(aim) != Bomb {
		return false
	}
	switch {
	case dk.GetBrandType(d) != Bomb && dk.GetBrandType(aim) == Bomb:
		return true
	case dk.GetBrandType(d) == Bomb && dk.GetBrandType(aim) == Bomb: //都是bomb直接比大小即可
		return AllComparison(d, aim)
	default:
		if len(d.Bd) != len(aim.Bd) {
			return false
		}
		return AllComparison(d, aim)
	}
}

//GetBrandType 类型反馈
//必须先调用SortOrderAsc排序
func (dk *DoubleBuckle) GetBrandType(d DeckOfCards) string {
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

//JudgeStraight 判断是顺子,顺子里不能有2，小王和大王
func JudgeStraight(str DeckOfCards) bool {
	if len(str.Bd) < 5 {
		return false
	}
	for i := 0; i < len(str.Bd)-1; i++ {
		switch str.Bd[i+1].Code {
		case 15:
			return false
		case 16:
			return false
		case 17:
			return false
		}
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
		switch str.Bd[i+1].Code {
		case 15:
			return false
		case 16:
			return false
		case 17:
			return false
		}
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
		switch str.Bd[i+1].Code {
		case 15:
			return false
		case 16:
			return false
		case 17:
			return false
		}
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

//ResponseOnMessage 接收到执行的逻辑，在newhub定义中赋值给方法类型
func ResponseOnMessage(data []byte, hub *ws.Hub) error {
	msg := ws.Message{}
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}
	res, ok := msg.Data.(DeckOfCards)
	if !ok {
		return errors.New("data type have error")
	}
	if dbuck.BrandComparison(lastBrands, res) {
		//减去用户对应brand
		allUserBrands.UserBrands[msg.User] = deleteBrand(allUserBrands.UserBrands[msg.User], res)
	}
	sendMsg, err := json.Marshal(allUserBrands)
	if err != nil {
		return err
	}
	hub.Broadcast <- ws.Message{
		User:     msg.User,
		Data:     sendMsg,
		DateTime: msg.DateTime,
	}
	return nil
}

//删除第一个参数中，第二个参数的内容
func deleteBrand(divisor, dividend DeckOfCards) DeckOfCards {
	for _, v := range dividend.Bd {
		for i, val := range divisor.Bd {
			if v.Code == val.Code {
				if i == 0 {
					divisor.Bd = divisor.Bd[1:]
					continue
				}
				if i == len(divisor.Bd)-1 {
					divisor.Bd = divisor.Bd[:len(divisor.Bd)-1]
					continue
				}
				divisor.Bd = append(divisor.Bd[:i], divisor.Bd[i+1:]...)
			}
		}
	}
	return divisor
}

func registerAndStart(user string) {
	allUserBrands.UserBrands[user] = DeckOfCards{}
	if len(allUserBrands.UserBrands) == 4 {
		StartBrandGame()
	}
}

//StartBrandGame 起始
func StartBrandGame() {
	totalBrands := dbuck.LicensingCode()
	log.Println(totalBrands)
}

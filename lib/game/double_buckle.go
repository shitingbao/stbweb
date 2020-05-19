package game

import (
	"encoding/json"
	"errors"
	"stbweb/core"
	"stbweb/lib/ws"
	"time"
)

//DoubleBuckle 双
type DoubleBuckle struct{}

//AllUserBrands 保存四个用户对应数据
type AllUserBrands struct {
	UserBrands map[string]DeckOfCards
}

//下发给用户的数据包
//下一个出牌的玩家
//数据包
//出的牌
//出牌的人
//UserShowStatues页面显示,里面有该用户的，就显示不出（其他两种情况，出牌的就匹配出牌人即可，两个都没的说明还没轮到，就什么都不显示）
//是否成功提交
type sendBrandsData struct {
	NextUser        string
	ShowDataUser    string
	Data            AllUserBrands
	ShowData        DeckOfCards
	UserShowStatues map[string]bool
	Success         bool
}

var (
	dbuck          DoubleBuckle
	lastBrands     DeckOfCards   //保存上一次的数据
	lastBrandsUser string        //上一次是谁出的
	allUserBrands  AllUserBrands //保存四个人的所有数据
	brandsUser     []string      //用户名以及顺序，不同的位置提前调整，保存的是已经是合理的顺序
	showStart      map[string]bool
)

//LicensingCode 给定所有牌,10以后，11，12，13，14，15，16增量代表J，Q，K，A，2,小王，大王，其他逻辑需要重定义玩法的话使用Ox
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
//这里接受到出的牌，删除后，把结果传递回去
//中间要区分第一次出，不出和压上
//先验证是否符合出的逻辑,错误信息在私人通道反馈，错误信息只包含一个false
//CodeErr这个错误是出牌类型不符合
func ResponseOnMessage(data []byte, hub *ws.Hub) error {
	msg := ws.Message{}
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}
	res, ok := msg.Data.(DeckOfCards)
	if !ok {
		return errors.New("data type have error")
	}
	if dbuck.GetBrandType(res) == CodeErr {
		hubUserSend(msg.User, msg.DateTime, hub)
		return nil
	}
	if len(res.Bd) == 0 { //res.Bd内容为空就是不出，如果这次不出的下一个用户，是上一次出牌的人，那就是一轮不要
		showStart[msg.User] = true
		if getNextUser(msg.User) == lastBrandsUser { //如果这个要不起的人是最后一个，那就是过了一轮，把上一次保存的清空，开始新的一轮
			lastBrands, lastBrandsUser, showStart = DeckOfCards{}, "", make(map[string]bool)
		}
		hubSend(msg.User, lastBrandsUser, allUserBrands, res, showStart, true, msg.DateTime, hub)
		return nil
	}
	if len(lastBrands.Bd) == 0 || dbuck.BrandComparison(lastBrands, res) { ////上一次出牌长度为0说明是第一次出，或者压上，出牌成功，减去用户对应brand，结果给所有人反馈
		allUserBrands.UserBrands[msg.User] = deleteBrand(allUserBrands.UserBrands[msg.User], res)
		lastBrands, lastBrandsUser, showStart = res, msg.User, make(map[string]bool) //showStart出牌后也要清空，因为页面上只需要显示出的牌即可
		hubSend(msg.User, lastBrandsUser, allUserBrands, res, showStart, true, msg.DateTime, hub)
		return nil
	}
	hubUserSend(msg.User, msg.DateTime, hub) //执行到这里说明出的牌不比上一次的大，不能这么出，也是给私人通道发送false信息
	return nil
}

//公共信息展示
func hubSend(user, lbs string, ads AllUserBrands, res DeckOfCards, st map[string]bool, isSuccess bool, tm time.Time, hub *ws.Hub) {
	result := sendBrandsData{
		NextUser:        getNextUser(user),
		ShowDataUser:    lbs,
		Data:            ads,
		ShowData:        res,
		UserShowStatues: st,
		Success:         isSuccess,
	}
	hub.Broadcast <- ws.Message{
		User:     user,
		Data:     result,
		DateTime: tm,
	}
}

//私人信息，内容为错误信息，出牌大小或者出牌的规程错误
func hubUserSend(user string, tm time.Time, hub *ws.Hub) {
	result := sendBrandsData{
		Success: false,
	}
	hub.BroadcastUser <- ws.Message{
		User:     user,
		Data:     result,
		DateTime: tm,
	}
}

func getNextUser(user string) string {
	nextUser := ""
	for i, v := range brandsUser {
		if v == user {
			if i == len(brandsUser)-1 {
				nextUser = brandsUser[0]
				break
			}
			nextUser = brandsUser[i+1]
			break
		}
	}
	return nextUser
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

//RegisterAndStart 注册人员，满四人开始
func RegisterAndStart(user string) {
	brandsUser = append(brandsUser, user)
	allUserBrands.UserBrands[user] = DeckOfCards{}
	if len(allUserBrands.UserBrands) == 4 {
		startBrandGame()
	}
}

//StartBrandGame 起始，使用map随机的特性，将原始总数据分发给四个用户
//这里的第一次数据给定一个第一个出的用户，后续待定
func startBrandGame() {
	totalBrands := dbuck.LicensingCode()
	// log.Println(totalBrands)
	vm := make(map[int]Brand)
	for i, v := range totalBrands.Bd {
		vm[i] = v
	}
	tl := 0
	for _, v := range vm {
		switch tl % 3 {
		case 0:
			setAllUserBrands(brandsUser[0], v)
		case 1:
			setAllUserBrands(brandsUser[1], v)
		case 2:
			setAllUserBrands(brandsUser[2], v)
		case 3:
			setAllUserBrands(brandsUser[3], v)
		}
		tl++
	}
	result := sendBrandsData{
		NextUser: brandsUser[0],
		Data:     allUserBrands,
		ShowData: DeckOfCards{},
		Success:  true,
	}
	core.CardHun.Broadcast <- ws.Message{
		User:     brandsUser[0],
		Data:     result,
		DateTime: time.Now(),
	}
}

//给用户发牌
func setAllUserBrands(user string, brd Brand) {
	md := allUserBrands.UserBrands[user].Bd
	md = append(md, brd)
	dcs := DeckOfCards{
		Bd: md,
	}
	allUserBrands.UserBrands[user] = dcs
}

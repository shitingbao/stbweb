package chatroom

import (
	"net/http"
	"stbweb/core"

	"github.com/pborman/uuid"
)

type chat struct{}
type chatCreateData struct {
	RoomName string
	NumTotle int
	RoomType string
	Common   string
}

func init() {
	core.RegisterFun("chat", new(chat), true)
}

//Post 业务处理,post请求的例子
func (ap *chat) Post(arge *core.ElementHandleArgs) {
	if arge.APIInterceptionPost("create", new(chatCreateData), createRoom) {
		return
	}
}

//新建一个room对象
func createRoom(param interface{}, p *core.ElementHandleArgs) error {
	pm := param.(*chatCreateData)

	roomID := uuid.NewUUID().String()
	room := roomPool.Get().(chatRoom)

	rClient, err := newRoomClient(roomID)
	if err != nil {
		return err
	}
	room.RoomID = roomID
	room.RoomName = pm.RoomName
	room.HostName = p.Usr
	room.NumTotle = pm.NumTotle
	room.Num = 1
	room.RoomType = pm.RoomType
	room.Common = pm.Common
	room.roomClient = rClient

	saveRoom() //存入mongodb
	openTCP()  //打开一个tcp连接
	core.SendJSON(p.Res, http.StatusOK, map[string]bool{"success": true})
	return nil
}

func saveRoom() {}
func openTCP()  {}

package apis

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"zinx-lwh/mmo_game_zinx/core"
	"zinx-lwh/mmo_game_zinx/pb"
	"zinx-lwh/ziface"
	"zinx-lwh/znet"
)

type WorldChatApi struct {
	znet.BaseRouter
}

func (wc *WorldChatApi) Handle(request ziface.IRequest)  {
	fmt.Println("注册路由")
	//1.解析客户端传进来的proto协议
	proto_msg := &pb.Talk{}
	err := proto.Unmarshal(request.GetData(),proto_msg)
	if err != nil {
		fmt.Println("proto unmarshal error",err)
		return
	}
	//2.当前的聊天数据是属于哪个玩家发送的
	pid ,err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty error",err)
	}
	//3.根据pid得到对应的player对象
	player := core.WorldMgrObj.Players[pid.(int32)]
	//4.将这个消息广播给其他全部在线玩家
	player.Talk(proto_msg.Content)
}
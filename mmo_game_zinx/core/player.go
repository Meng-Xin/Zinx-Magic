package core

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"sync"
	"zinx-lwh/mmo_game_zinx/pb"
	"zinx-lwh/ziface"
)

// Player 玩家对象
type Player struct {
	Pid  int32              //玩家ID
	Conn ziface.IConnection //为当前玩家建立的连接（用于和客户端建立连接）
	X    float32            //平面X 坐标
	Y    float32            //高度
	Z    float32            //平面Y 坐标
	V    float32            //玩家旋转角度 0-360角度

}

/*
	Player ID 生成器
*/

var PidGen int32 = 1  //用来生成玩家ID的计数器
var IdLock sync.Mutex //保护PidGen的Mutex

// NewPlayer  创建一个玩家的方法
func NewPlayer(conn ziface.IConnection) *Player {
	//生成一个玩家ID
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()

	//实例化一个玩家对象
	player := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)), //随机在160坐标点，基于X轴进行偏移
		Y:    0,
		Z:    float32(140 + rand.Intn(20)), //随机在140坐标点，基于Y轴进行偏移
		V:    0,                            //玩家默认角度0度
	}
	return player
}

/*
	提供一个发送给客户端消息的方法
	主要是将pb的protobuf数据序列化之后，在调用zinx的SendMsg方法
*/

func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	//1.将待发送数据进行序列化proto处理， Message结构体序列化，转成二进制
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err :", err)
		return
	}

	//2.先判断当前连接是否断开
	if p.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}
	//3.将序列化后的二进制文件，通过zinx框架的sendMsg将数据发送给客户端
	err = p.Conn.SendMsg(msgId, msg)
	if err != nil {
		fmt.Println("Player sendMsg err", err)
		return
	}
}

// SyncPid 告知客户端玩家Pid，同步已经生成的玩家ID给客户端
func (p *Player) SyncPid() {
	//组件MsgID:1 的proto数据
	proto_msg := &pb.SyncPid{
		Pid: p.Pid,
	}
	//将消息发送给客户端
	p.SendMsg(1, proto_msg)
}

// BroadCastStarPosition 广播玩家自己的出生地点
func (p *Player) BroadCastStarPosition() {
	//组建MsgID:200 的proto数据
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2, //TP2 代表广播的是位置坐标
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	//将消息发送给客户端
	p.SendMsg(200, proto_msg)
}

// Talk 玩家广播世界聊天消息
func (p *Player) Talk (content string)  {
	//1.组件MsgID:200 proto数据
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp: 1,		//1-代表聊天
		Data:&pb.BroadCast_Content{
			Content: content,
		},
	}
	//2.得到当前世界中的在线玩家
	players := WorldMgrObj.GetAllPlayers()
	//3.向所有玩家（包括自己）发送MsgID:200消息s
	for _,player := range players{
		//player分别给对用的客户端发送消息
		player.SendMsg(200,proto_msg)
	}

}
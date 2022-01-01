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
func (p *Player) Talk(content string) {
	//1.组件MsgID:200 proto数据
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  1, //1-代表聊天
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}
	//2.得到当前世界中的在线玩家
	players := WorldMgrObj.GetAllPlayers()
	//3.向所有玩家（包括自己）发送MsgID:200消息s
	for _, player := range players {
		//player分别给对用的客户端发送消息
		player.SendMsg(200, proto_msg)
	}

}

// SyncSurrounding 同步玩家上线的位置信息
func (p *Player) SyncSurrounding() {
	//1。获取当前玩家周围的玩家有哪些（九宫格）
	pids := WorldMgrObj.AoiMgr.GetPidByPos(p.X, p.Z)
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}
	//2。将当前玩家的位置信息通过MsgID:200 发送给周围的玩家（让其他玩家看到自己）
	//2.1组建MsgID:200的proto数据
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2, //2- 代表广播坐标
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	//2.2 全部周围的玩家都向格子的客户端发送200消息，proto_msg
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}
	//3。将周围的全部玩家的位置信息发送给当前的玩家客户端 MsgID:202
	//3.1 制作MsgID:202 proto数据
	//3.1.1 制作pb.Player Slice
	players_proto_msg := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		//制作一个Player Message
		p := &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		players_proto_msg = append(players_proto_msg, p)
	}
	//3.1.2 封装SyncPlayer protobuf数据
	Syncplayers_proto_msg := &pb.SyncPlayers{
		Ps: players_proto_msg[:],
	}
	//3.2 将组建好的数据发送给当前游戏玩家的客户端
	p.SendMsg(202, Syncplayers_proto_msg)
}


// UpdatePos 广播当前玩家的位置信息
func (p *Player) UpdatePos(x,y,z,v float32)  {
	//1.更新当前玩家的player对象坐标
	p.X,p.Y,p.Z,p.V = x,y,z,v
	//2.组建广播proto协议，MsgID:200 tp-4
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp: 4,	//4-移动之后的坐标信息
		Data: &pb.BroadCast_P{
			P:&pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	//3.获取当前玩家的周边玩家AOI九宫格之内的玩家
	players := p.GetSurroundingPlayers()
	//4.依次给灭个玩家对应的客户端发送当前玩家的位置更新信息
	for _,player := range players{
		player.SendMsg(200,proto_msg)
	}
}


// GetSurroundingPlayers 获取当前玩家的周边玩家AOI九宫格之内的玩家
func (p *Player) GetSurroundingPlayers() []*Player {
	//得到当前AOI九宫格内的所有玩家ID
	pids := WorldMgrObj.AoiMgr.GetPidByPos(p.X,p.Z)

	//将所有的pid对应的Player放到Players切片中
	players := make([]*Player,0,len(pids))
	for _,pid := range pids{
		players = append(players,WorldMgrObj.GetPlayerByPid(int32(pid)))
	}
	return players
}

//
func (p *Player) Offline ()  {
	//得到当前玩家周围的九宫格内的所有玩家
	players := p.GetSurroundingPlayers()
	//给周围玩家广播MsgID:201消息
	proto_msg := &pb.SyncPid{
		Pid: p.Pid,
	}
	for _,player := range players{
		player.SendMsg(201,proto_msg)
	}

	//WorldMgrObj.AoiMgr.RemoveFromGridByPos(int(p.Pid),p.X,p.Z)
	WorldMgrObj.RemovePlayerByPid(p.Pid)
}
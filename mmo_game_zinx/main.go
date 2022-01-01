package main

import (
	"fmt"
	"zinx-lwh/mmo_game_zinx/apis"
	"zinx-lwh/mmo_game_zinx/core"
	"zinx-lwh/ziface"
	"zinx-lwh/znet"
)

// OnConnectionAdd 当前客户端建立连接之后的hook函数
func OnConnectionAdd(conn ziface.IConnection)  {
	fmt.Println("调用注册方法")
	//创建一个Player对象
	player := core.NewPlayer(conn)
	//给客户端发送MsgID:1的消息:同步当前Player的ID 给客户端
	player.SyncPid()
	//给客户端发送MsgID：200的消息，同步当前Player 的位置给客户端
	player.BroadCastStarPosition()

	//将当前新上线的玩家添加到WorldManager中
	core.WorldMgrObj.AddPlayer(player)

	//将该链接绑定一个Pid 玩家id属性
	conn.SetProperty("pid",player.Pid)

	//同步周边玩家，告知他们当前玩家已经上线，广播当前玩家位置信息
	player.SyncSurrounding()

	fmt.Println("===> Player pid=",player.Pid,"is arrived! <===")
}

// OnConnectionLost 给当前连接断开之前出发的Hook钩子函数
func OnConnectionLost(conn ziface.IConnection)  {
	// 通过连接属性得到当前连接所绑定的pid
	pid,_ := conn.GetProperty("pid")
	//得到当前玩家周围的九宫格内的所有玩家
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	//触发玩家下线的业务
	player.Offline()
	//给周围玩家发送
	fmt.Println("===>Player pid=",pid,"offline<===")
}

func main() {
	//创建Zinx server句柄
	s := znet.NewServer("Zinx Game V0.10")
	//连接创建和销毁HOOK钩子函数
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)
	//注册一些路由业务
	s.AddRouter(2,&apis.WorldChatApi{})
	s.AddRouter(3,&apis.MoveApi{})
	//启动服务
	s.Server()

}
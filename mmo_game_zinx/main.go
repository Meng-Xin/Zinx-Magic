package main

import (
	"fmt"
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

	fmt.Println("===> Player pid=",player.Pid,"is arrived! <===")
}



func main() {
	//创建Zinx server句柄
	s := znet.NewServer("Zinx Game V0.10")
	//连接创建和销毁HOOK钩子函数
	s.SetOnConnStart(OnConnectionAdd)
	//注册一些路由业务

	//启动服务
	s.Server()

}
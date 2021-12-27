package main

import (
	"fmt"
	"zinx-lwh/ziface"
	"zinx-lwh/znet"
)

/*
	基于Zinx框架来开发的 服务器应用程序
*/
//ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Handle 测试
func (b *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")
	fmt.Println("recv from client:msgId:",request.GetMsgID(),"msgData:",string(request.GetData()))
	err := request.GetConnection().SendMsg(200,[]byte("ping...ping...ping..."))
	if err != nil {
		return
	}
}

// HelloZinxRouter  自定义路由测试
type HelloZinxRouter struct {
	znet.BaseRouter
}
func (b *HelloZinxRouter) Handle (request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle")
	fmt.Println("recv from client:msgId:",request.GetMsgID(),"msgData:",string(request.GetData()))
	err := request.GetConnection().SendMsg(201,[]byte("Hello Welcome to  Zinx "))
	if err != nil {
		return
	}
}


func main() {
	//1、创建一个Server句柄，使用Zinx的Api
	s := znet.NewServer("[zinx V0.7]")

	//2. 给当前zinx框架添加一个自定义router
	s.AddRouter(0,&PingRouter{})
	s.AddRouter(1,&HelloZinxRouter{})

	//3、启动server
	s.Server()
}

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
// PreHandle 测试
func (b *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle")
	_,err := request.GetConnection().GetTCPConnection().Write([]byte("before ping..."))
	if err != nil {
		fmt.Println("call back before ping... error:",err)
	}
}
// Handle 测试
func (b *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")
	_,err := request.GetConnection().GetTCPConnection().Write([]byte("Handle ping..."))
	if err != nil {
		fmt.Println("call back ping... ping...ping... error:",err)
	}
}
// PostHandle 测试
func (b *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle")
	_,err := request.GetConnection().GetTCPConnection().Write([]byte("PostHandle ping..."))
	if err != nil {
		fmt.Println("call back after ping... error:",err)
	}
}


func main() {
	//1、创建一个Server句柄，使用Zinx的Api
	s := znet.NewServer("[zinx V0.3]")

	//2. 给当前zinx框架添加一个自定义router
	s.AddRouter(&PingRouter{})

	//3、启动server
	s.Server()
}

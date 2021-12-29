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
// DoConnectionBegin 创建链接之后执行的钩子函数
func DoConnectionBegin(conn ziface.IConnection)  {
	fmt.Println("===>DoConnectionBegin is Called! ...")
	err := conn.SendMsg(202, []byte("DoConnection BEGIN"))
	if err != nil {
		fmt.Println("Call DoConnectionBegin error:",err)
	}
	//给当前连接设置属性
	fmt.Println("Set Conn Property ...")
	conn.SetProperty("Name","李文豪-Meng-Xin")
	conn.SetProperty("GitHub","https://github.com/meng-xin")
	conn.SetProperty("Blog","https://www.jianshu.com/u/be6558c1e3b6")
}
// DoConnectionBefore 断开连接之前需要执行的钩子函数
func DoConnectionBefore(conn ziface.IConnection)  {
	fmt.Println("===> DoConnectionBefore is Called ...")
	fmt.Println("conn ID",conn.GetConnID(),"is Lost")

	//获取连接属性
	if name,err := conn.GetProperty("Name");err ==nil{
		fmt.Println("Name:",name)
	}
	if github,err := conn.GetProperty("GitHub");err == nil {
		fmt.Println("GitHub:",github)
	}
	if blog,err := conn.GetProperty("Blog"); err == nil {
		fmt.Println("Blog:",blog)
	}
}
func main() {
	//1、创建一个Server句柄，使用Zinx的Api
	s := znet.NewServer("[zinx V0.7]")
	//2、注册链接Hook 钩子函数回调
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionBefore)
	//3、 给当前zinx框架添加一个自定义router
	s.AddRouter(0,&PingRouter{})
	s.AddRouter(1,&HelloZinxRouter{})

	//4、启动server
	s.Server()
}

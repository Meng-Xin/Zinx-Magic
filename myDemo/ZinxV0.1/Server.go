package main

import "zinx-lwh/znet"

// main 基于Zinx框架来开发的 服务器应用程序
func main() {
	//1、创建一个Server句柄，使用Zinx的Api
	s := znet.NewServer("[zinx V0.1]")
	//2、启动server
	s.Server()
}

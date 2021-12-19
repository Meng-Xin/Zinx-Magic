package znet

import (
	"fmt"
	"net"
	"zinx-lwh/ziface"
)

type Server struct {
	//服务器的名称
	Name string
	//服务器绑定的ip版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int
}

func (s *Server)Start()  {
	fmt.Printf("[Strat] Server Listenner at IP :%s, Port %d, is starting\n",s.IP,s.Port)
	go func() {
		//1、获取一个TCP的Addr
		addr,err := net.ResolveTCPAddr(s.IPVersion,fmt.Sprintf("%s:%d",s.IP,s.Port))
		if err != nil {
			fmt.Println("resolve tcp addt error : ",err)
			return
		}
		//2、监听服务器地址
		listenner,err := net.ListenTCP(s.IPVersion,addr)
		if err != nil {
			fmt.Println("Listen",s.IPVersion,"err:",err)
			return
		}
		fmt.Println("start Zinx server succ,",s.Name," succ Listening...")
		//3、阻塞等待客户端连接，处理客户端链接的业务（读写）
		for  {
			//如果有客户端连接过来，阻塞会返回
			conn,err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err",err)
				continue
			}
			//如果与客户端建立连接，做一些业务，做一个最基本的最大512bit字节长度的回显业务
			go func() {
				for  {
					buf := make([]byte,512)
					cnt,err := conn.Read(buf)
					if err != nil {
						fmt.Println("recv buf err",err)
						continue
					}
					fmt.Println("recv client buf %s, cnt %d",buf,cnt)
					//回显功能
					if _,err := conn.Write(buf[:cnt]);err !=nil {
						fmt.Println("write back buff err:",err)
						continue
					}
				}
			}()
		}
	}()
}
func (s *Server)Server()  {
	s.Start()
	//TODO 做一些启动服务器之后的额外业务

	// 阻塞状态
	select {

	}
}
func (s *Server)Stop()  {
	//TODO 将一些服务器的资源、状态或者一些已经开辟的链接信息，进行停止或者回收

}

// 初始化Server模块的方法

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name : name,
		IPVersion: "tcp4",
		IP: "0.0.0.0",
		Port: 8999,
	}
	return s
}
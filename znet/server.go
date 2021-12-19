package znet

import (
	"errors"
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
//TODO CallBackToClient 定义当前客户端链接的所绑定的handle api(目前这个handle是写死的，以后优化应该由用户自定义handle方法)
func CallBackToClient (conn *net.TCPConn,data []byte,cnt int) error{
	//回显的业务
	fmt.Println("[Conn Handle] CallBackToClient...")
	if _,err := conn.Write(data[:cnt]); err != nil{
		fmt.Println("write back buf err:",err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

// Start 启动服务器
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
		var cid uint32
		cid = 0
		//3、阻塞等待客户端连接，处理客户端链接的业务（读写）
		for  {
			//如果有客户端连接过来，阻塞会返回
			conn,err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err",err)
				continue
			}
			//处理新连接的业务方法和conn进行绑定 得到我们的链接模块
			dealConn := NewConnection(conn,cid,CallBackToClient)
			cid++

			//启动当前连接业务的处理方法
			go dealConn.Start()
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
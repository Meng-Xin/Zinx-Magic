package znet

import (
	"fmt"
	"net"
	"zinx-lwh/utils"
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
	//当前server 的消息管理模块 ，用例绑定MsgID和对应的处理业务API方法
	MsgHandler ziface.IMsgHandle

}
// AddRouter 路由功能：给当前的服务注册一个路由方法，供客户端的链接处理使用
func (s *Server) AddRouter (msgID uint32,router ziface.IRouter)  {
	s.MsgHandler.AddRouter(msgID,router)
	fmt.Println("Add Router Succ!!")
}



// Start 启动服务器
func (s *Server)Start()  {
	fmt.Printf("[Zinx] Server Name : %s,Listenner at IP : %s,Port:%d is strating",
		utils.GlobalObject.Name,utils.GlobalObject.Host,utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version %s, MaxConn:%d,MaxPackeetSize:%d\n",
		utils.GlobalObject.Version,utils.GlobalObject.MaxConn,utils.GlobalObject.MaxPackageSize)
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
			dealConn := NewConnection(conn,cid,s.MsgHandler)
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
		Name : utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP: utils.GlobalObject.Host,
		Port: utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
	}
	return s
}
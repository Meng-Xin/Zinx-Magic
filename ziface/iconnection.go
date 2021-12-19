package ziface

import "net"

// IConnection 定义链接模块的抽象层
type  IConnection interface {
	// Start 启动链接 让当前的连接准备开始工作
	Start()

	// Stop  停止链接 结束当前连接的工作
	Stop()

	//GetTCPConnection 获取当前连接绑定的Socket conn
	GetTCPConnection() *net.TCPConn

	// GetConnID 获取当前连接模块的连接ID
	GetConnID() uint32

	// RemoteAddr 获取远程客户端的TCP状态 IP PORT
	RemoteAddr() net.Addr

	// Send 发送数据，将数据发送给远程客户端
	Send(data []byte)error
}

// HandleFunc 定义一个处理链接业务的方法
type HandleFunc func(*net.TCPConn,[]byte,int) error
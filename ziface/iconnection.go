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

	// RemoteAddr 获取远程客户端的TCP状态 IP port
	RemoteAddr() net.Addr

	// SendMsg 发送数据，将数据发送给远程客户端
	SendMsg(msgId uint32, data []byte) error

	// SetProperty 设置连接属性
	SetProperty(key string,value interface{})
	// GetProperty 获取连接属性
	GetProperty(key string) (interface{},error)
	// RemoveProperty 移除连接属性
	RemoveProperty(key string)
}

// HandleFunc 定义一个处理链接业务的方法
type HandleFunc func(*net.TCPConn,[]byte,int) error
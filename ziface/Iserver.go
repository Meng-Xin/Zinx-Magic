package ziface

type IServer interface {
	//Start 启动服务器
	Start()
	// Stop 停止服务器
	Stop()
	// Server 运行服务器
	Server()
	// AddRouter 路由功能：给当前的服务注册一个路由方法，供客户端的链接处理使用
	AddRouter(msgID uint32, router IRouter)
	// GetConnMgr 获取当前server 的链接管理器
	GetConnMgr() IConnManager
	// SetOnConnStart 注册OnConnStart 钩子函数的方法
	SetOnConnStart(func(conn IConnection))
	// SetOnConnStop 注册OnConnStop 钩子函数的方法
	SetOnConnStop(func(conn IConnection))
	// CallOnConnStart 调用OnConnStart 钩子函数的方法
	CallOnConnStart(connection IConnection)
	// CallOnConnStop 调用OnConnStop 钩子函数的方法
	CallOnConnStop(connection IConnection)
}

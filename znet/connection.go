package znet

import (
	"fmt"
	"net"
	"zinx-lwh/ziface"
)

// Connection 链接模块
type Connection struct {
	//当前链接的Socket TCP 套接字
	Conn *net.TCPConn
	//链接的ID
	ConnID uint32
	//当前的连接状态
	isClosed bool
	//当前连接所绑定的处理业务的方法API
	handleAPI ziface.HandleFunc
	//告知当前连接已经退出的/停止 Channel
	ExitChan chan bool
}

func (c *Connection)StartReader()  {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID:",c.ConnID,"Reader is exit remote addr is",c.RemoteAddr().String())
	defer c.Stop()

	for  {
		//读取客户端的数据到buf中，最大512字节
		buf := make([]byte,512)
		cnt,err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err:",err)
			continue
		}
		//调用当前链接所绑定的HandleAPI
		if err := c.handleAPI(c.Conn,buf,cnt);err!=nil{
			fmt.Println("ConnID ",c.ConnID,"handleAPI is err:",err)
			break
		}
	}
}
// Start 启动链接 让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start() ...ConnID=",c.ConnID)
	//启动从当前链接的读取数据业务
	go c.StartReader()
	//TODO 启动从当前链接写数据的业务
}
// Stop  停止链接 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID=",c.ConnID)
	//如果当前连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	//关闭socket
	c.Conn.Close()
	//回收资源
	close(c.ExitChan)
}

//GetTCPConnection 获取当前连接绑定的Socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前连接模块的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获取远程客户端的TCP状态 IP PORT
func (c*Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// Send 发送数据，将数据发送给远程客户端
func (c *Connection) Send(data []byte) error {
	panic("implement me")
}

// 初始化链接模块的方法
func NewConnection(conn *net.TCPConn,connID uint32,callback_api ziface.HandleFunc) *Connection  {
	c := &Connection{
		conn,
		connID,
		false,
		callback_api,
		make(chan bool ,1),
	}
	return c
}
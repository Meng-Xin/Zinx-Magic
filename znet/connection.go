package znet

import (
	"errors"
	"fmt"
	"io"
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

	//告知当前连接已经退出的/停止 Channel
	ExitChan chan bool
	//该链接处理的方法Router
	Router ziface.IRouter
}

func (c *Connection)StartReader()  {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID:",c.ConnID,"Reader is exit remote addr is",c.RemoteAddr().String())
	defer c.Stop()

	for  {
		//读取客户端的数据到buf中，最大为用户配置大小
		//buf := make([]byte,utils.GlobalObject.MaxPackageSize)
		//_,err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("recv buf err:",err)
		//	continue
		//}
		//创建一个拆包对象
		dp := NewDataPack()

		//读取客户端的Msg Head 8 个字节
		headData := make([]byte,dp.GetHeadLen())
		//拆包，得到MsgID 和 msgDataLen 放在msg消息中
		if _,err := io.ReadFull(c.GetTCPConnection(),headData);err != nil{
			fmt.Println("read msg head error:",err)
			break
		}
		//创建存储msg对象,进行拆包
		msg,err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error:",err)
			break
		}
		//根据dataLen 再次读取Data，放在msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte,msg.GetDataLen())
			if _,err := io.ReadFull(c.GetTCPConnection(),data);err != nil{
				fmt.Println("read msg data error",err)
				break
			}
		}
		msg.SetData(data)

		//得到当前conn 链接的request请求数据
		req:=Request{
			conn: c,
			msg: msg,
		}
		//执行注册的路由方法
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
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

// SendMsg 提供一个SendMsg方法，将要发送给客户端的数据，线进行封包，再发送
func (c *Connection)SendMsg(msgId uint32, data []byte) error  {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}
	//将data 进行封包,MsgDataLen | MsgId |MsgData
	dp := NewDataPack()
	//MsgDataLen | MsgId | MsgData
	binaryMsg,err := dp.Pack(NewMsgPackage(msgId,data))
	if err != nil {
		fmt.Println("pack error msg Id:",msgId)
		return errors.New("Pack error msg")
	}
	//将数据 发送给客户端
	if _,err := c.Conn.Write(binaryMsg);err != nil{
		fmt.Println("Write msg id",msgId,"error",err)
		return errors.New("conn Write err")
	}
	return nil
}

// NewConnection 初始化链接模块的方法
func NewConnection(conn *net.TCPConn,connID uint32,router ziface.IRouter) *Connection  {
	c := &Connection{
		Conn:conn,
		ConnID: connID,
		Router: router,
		isClosed: false,
		ExitChan: make(chan bool ,1),
	}
	return c
}
package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx-lwh/utils"
	"zinx-lwh/ziface"
)

// Connection 链接模块
type Connection struct {
	//当前Conn隶属于那个Server
	TcpServer ziface.IServer
	//当前链接的Socket TCP 套接字
	Conn *net.TCPConn
	//链接的ID
	ConnID uint32
	//当前的连接状态
	isClosed bool

	//告知当前连接已经退出的/停止 Channel
	ExitChan chan bool
	//无缓冲的管道，用于读、写Goroutine之间的消息通信
	msgChan chan []byte

	//消息的管理MsgID 和对应的处理业务API关系
	MsgHandler ziface.IMsgHandle

	//链接属性集合
	property map[string]interface{}
	//保护链接属性的锁
	propertyLock sync.RWMutex
}

// NewConnection 初始化链接模块的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer: server,
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
		property: make(map[string]interface{}), //TODO 这里我自己添加的
	}
	//将conn 加入到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running...]")
	defer fmt.Println("connID:", c.ConnID, "[Reader is exit!], remote addr is", c.RemoteAddr().String())
	defer c.Stop()

	for {
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
		headData := make([]byte, dp.GetHeadLen())
		//拆包，得到MsgID 和 msgDataLen 放在msg消息中
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error:", err)
			break
		}
		//创建存储msg对象,进行拆包
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error:", err)
			break
		}
		//根据dataLen 再次读取Data，放在msg.Data中 TODO 感觉有优化空间，这个data写法
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error", err)
				break
			}
		}
		msg.SetData(data)

		//得到当前conn 链接的request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经开启了工作池机制，将消息发送给workerPool 处理即可
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			//从路由中，找到注册绑定的Conn对应的router调用
			//根据绑定好的MsgID 找到对应处理api业务 执行
			go c.MsgHandler.DoMsgHandler(&req)
		}

	}
}

/*
	写消息的Goroutine，专门发送给客户端消息的模块
*/
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")
	//不断的阻塞的等待channel的消息，把消息写给客户端
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error,", err)
				return
			}
		case <-c.ExitChan:
			//代表Reader一经推出，此时Writer也需要退出
			return
		}
	}
}

// Start 启动链接 让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start() ...ConnID=", c.ConnID)
	//启动从当前链接的读取数据业务
	go c.StartReader()
	// 启动从当前链接写数据的业务
	go c.StartWriter()

	//按照开发者传递进来的 创建链接之后需要调用的处理业务，执行对应的Hook方法
	c.TcpServer.CallOnConnStart(c)
}

// Stop  停止链接 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID=", c.ConnID)
	//如果当前连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	//调用开发者注册的，销毁链接之前，需要执行对用的Hook方法
	c.TcpServer.CallOnConnStop(c)
	//关闭socket
	c.Conn.Close()
	//关闭Writer 业务，告知Writer 关闭
	c.ExitChan <- true
	//将当前连接从ConnMgr 中摘除
	c.TcpServer.GetConnMgr().Remove(c)
	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
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
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SendMsg 提供一个SendMsg方法，将要发送给客户端的数据，线进行封包，再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}
	//将data 进行封包,MsgDataLen | MsgId |MsgData
	dp := NewDataPack()
	//MsgDataLen | MsgId | MsgData
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("pack error msg Id:", msgId)
		return errors.New("Pack error msg")
	}
	//将数据 发送给 管道，通过管道 Writer 给客户端
	c.msgChan <- binaryMsg

	return nil
}

// SetProperty 设置连接属性
func (c *Connection) SetProperty(key string,value interface{})  {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

// GetProperty 获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if val,ok := c.property[key];ok {
		return val,nil
	}
	return nil,errors.New("no property found")
}

// RemoveProperty 移除连接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property,key)
}
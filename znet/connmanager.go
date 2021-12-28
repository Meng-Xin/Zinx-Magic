package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx-lwh/ziface"
)

/*
	连接管理模块
*/

type ConnManager struct {
	connections map[uint32]ziface.IConnection //管理的链接集合
	connLock    sync.RWMutex                  //保护链接集合的读写锁
}

// NewConnManager 创建当前链接实例方法

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// Add 添加链接
func (c *ConnManager) Add(conn ziface.IConnection) {
	//保护共享资源map ,加 写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	//将conn加入到ConnManager中
	c.connections[conn.GetConnID()] = conn
	fmt.Println("connID=",conn.GetConnID()," add to ConnManager SuccessFully: conn num:",c.Len())
}

// Remove 删除链接
func (c *ConnManager) Remove(conn ziface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	//删除conn链接信息
	delete(c.connections,conn.GetConnID())
	fmt.Println("connID=",conn.GetConnID()," remove from ConnManager SuccessFully: conn num:",c.Len())
}

// Get 根据ConnID获取链接
func (c *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	c.connLock.RLock()
	defer c.connLock.RUnlock()
	if conn,ok := c.connections[connID]; ok {
		//找到了
		return conn,nil
	}
	return nil,errors.New("connection not Found ")
}

// Len 得到当前链接总数
func (c *ConnManager) Len() int {
	return len(c.connections)
}

// ClearConn 清除并释放所有连接
func (c *ConnManager) ClearConn() {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	//删除conn并停止conn的工作
	for connID,conn := range c.connections{
		//停止
		conn.Stop()
		//删除
		delete(c.connections,connID)
	}
	fmt.Println("Clear All connections success! conn num = ",c.Len())
}

package znet

import "zinx-lwh/ziface"

type Request struct {
	//已经和客户端建立好的链接
	conn ziface.IConnection
	//客户端的数据
	msg ziface.IMessage
}
// GetConnection 得到链接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}
// GetData 客户端请求数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request)GetMsgID() uint32  {
	return r.msg.GetMsgId()
}

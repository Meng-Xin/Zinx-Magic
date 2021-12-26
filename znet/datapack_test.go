package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)
// 只是负责测试datapack拆包， 封包的单元测试
func TestDataPack(t *testing.T)  {
	// 模拟 | 服务器 |
	//1. 创建socketTCP
	listenner,err:=net.Listen("tcp","127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err :",err)
		return
	}
	//创建一个go承载，负责从客户端处理业务
	go func() {
		//2. 从客户端读取数据，拆包处理
		for  {
			conn,err :=listenner.Accept()
			if err != nil {
				fmt.Println("server accept error:",err)
				break
			}
			//通过协程 处理 conn链接，不断读取
			go func(conn net.Conn) {
				//处理客户端的请求
				// ----> 拆包的过程 <----
				//定义一个拆包对象
				dp := NewDataPack()
				for  {
					//1.第一次从conn读，把包的head读出来
					headData := make([]byte,dp.GetHeadLen())
					//通过io.ReadFull 通过io流读取，一次性读满数据。
					_,err := io.ReadFull(conn,headData)
					if err != nil {
						fmt.Println("read head err :",err)
						break
					}
					msgHead,err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err :",err)
						return
					}
					//判断 msg 是否有数据，有数据就进行第二次读取
					if msgHead.GetDataLen() > 0 {
						//2.第二次从conn读,把head中的dataLen读出来，在读取data内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte,msgHead.GetDataLen())

						//根据dataLen的长度再次从io流中读取
						_,err := io.ReadFull(conn,msg.Data)
						if err != nil {
							fmt.Println("server unpack err: ",err)
							return
						}
						//完整度的一个消息已经读取完毕
						fmt.Println("-->Recv MsgID",msg.Id,"DataLen:",msg.DataLen,"data：",string(msg.Data))
					}
				}

			}(conn)
		}
	}()


	// 模拟 | 客户端 |
	conn,err := net.Dial("tcp","127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err:",err)
		return
	}

	//创建一个封包对象 dp
	dp := NewDataPack()

	//模拟粘包过程，封装两个msg一同发送
	//封装第一个msg1包
	msg1 := &Message{
		Id:1,
		DataLen: 4,
		Data: []byte{'z','i','n','x'},
	}
	sendData1,err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err:",err)
		return
	}
	//封装第二个msg2包
	msg2 := &Message{
		Id:2,
		DataLen: 7,
		Data: []byte{'h','e','l','l','o','!','!'},
	}
	sendData2,err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg1 err:",err)
		return
	}
	//将两个包黏在一起
	sendData1 = append(sendData1,sendData2...)
	//一次性发送两个包
	conn.Write(sendData1)
	//客户端阻塞
	select {
	}
}

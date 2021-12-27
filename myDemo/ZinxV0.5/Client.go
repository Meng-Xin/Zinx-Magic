package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx-lwh/znet"
)

// main() 模拟客户端
func main()  {
	fmt.Println("client start...")
	//1.直接连接远程服务器。得到一个conn连接
	time.Sleep(1 *time.Second)
	conn,err := net.Dial("tcp","127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err:",err)
		return 
	}
	for  {
		//发送封包的message消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0,[]byte("ZinxV0.5 client Test Message!!!")))
		if err != nil {
			fmt.Println("Client Pack error:",err)
			return
		}
		//发送封包数据给服务端
		if _, err := conn.Write(binaryMsg); err != nil{
			fmt.Println("Client Write error:",err)
			return
		}
		//发送成功后，接收服务器的回复数据。 ：MsgID：1| ping...ping...
		
		//1.先读取流中的head部分，得到ID和dataLen
		binaryHead := make([]byte,dp.GetHeadLen())
		if _,err := io.ReadFull(conn,binaryHead);err != nil{
			fmt.Println("Client read head error:",err)
			break
		}
		//2.将二进制的binaryHead 进行unpack
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("Client unpakc error:",err)
			return
		}
		if msgHead.GetDataLen() > 0 {
			//2.让后根据 dataLen进行二次读取，得到MsgData
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte,msgHead.GetDataLen())
			if _,err := io.ReadFull(conn,msg.Data);err != nil{
				fmt.Println("Client read msg data error:",err)
				break
			}
			fmt.Println("--->Recv Server Msg:ID=",msg.GetMsgId(),"data=",string(msg.GetData()))
		}

		//cpu 阻塞
		time.Sleep(1*time.Second)
	}


}


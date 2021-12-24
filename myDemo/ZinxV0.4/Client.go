package main

import (
	"fmt"
	"net"
	"time"
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
		//2.连接调用write 写入数据
		_,err := conn.Write([]byte("Hello Zinx V0.2..."))
		if err != nil {
			fmt.Println("write conn err:",err)
			return
		}
		buf := make([]byte,512)
		cnt,err := conn.Read(buf)
		if err != nil{
			fmt.Println("read buf err:",err)
			return
		}
		fmt.Printf("server call back %s, cnt=%d\n",buf,cnt)
		//cpu 阻塞
		time.Sleep(1*time.Second)
	}


}


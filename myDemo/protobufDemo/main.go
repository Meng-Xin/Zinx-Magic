package main

import (
	"fmt"
	"zinx-lwh/myDemo/protobufDemo/pb"
	"github.com/golang/protobuf/proto"
)

func main() {
	person := &pb.Person{
		Name: "Li XiaoMing",
		Age: 18,
		Emails: []string{"997266754@qq.com","li997266754@gmail.com"},
		Phones: []*pb.PhoneNumber{
			&pb.PhoneNumber{
				Number: "14725836911",
				Type: pb.PhoneType_MOBILE,
			},
			&pb.PhoneNumber{
				Number: "12345678911",
				Type: pb.PhoneType_HOME,
			},
			&pb.PhoneNumber{
				Number: "15617242236",
				Type: pb.PhoneType_WORK,
			},
		},
	}

	//编码！
	//将 person 结构体（待发送数据 ）进行序列化处理---protobuf,得到一个二进制文件
	data ,err := proto.Marshal(person)
	if err != nil {
		fmt.Println("protobuf marshal error :",err)
	}
	//序列化完成的 data 就是需要通过网络传输的数据,对端需要按照Message Person格式进行解析

	//解码
	newPerson := &pb.Person{}
	err = proto.Unmarshal(data, newPerson)
	if err != nil {
		fmt.Println("protobuf Unmarshal error:",err)
	}
	fmt.Println("源数据:",person)
	fmt.Println("解码之后的数据:",newPerson)
}
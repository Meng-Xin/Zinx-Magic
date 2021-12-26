package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx-lwh/utils"
	"zinx-lwh/ziface"
)

// DataPack 封包，拆包具体模块
type DataPack struct {

}
// NewDataPack 结构体实例初始化方法
func NewDataPack() *DataPack  {
	return &DataPack{}
}

// GetHeadLen 获取包的头的长度方法
func (d *DataPack) GetHeadLen() uint32 {
	//DataLen uint32 (4字节) + ID uint32(4字节)
	return 8
}

// Pack 封包方法: 封装顺序-->|dataLen|msg|data|
func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲区
	dataBuff := bytes.NewBuffer([]byte{})
	//1.将dataLen 写进dataBuff中
	if err := binary.Write(dataBuff,binary.LittleEndian,msg.GetDataLen());err != nil {
		return nil,err
	}
	//2.将MsgId  写进dataBuff中
	if err := binary.Write(dataBuff,binary.LittleEndian,msg.GetMsgId());err != nil {
		return nil,err
	}
	//3.将data数据 写进dataBuff中
	if err := binary.Write(dataBuff,binary.LittleEndian,msg.GetData());err != nil {
		return nil,err
	}
	return dataBuff.Bytes(),nil
}

// Unpack 拆包方法
func (d *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	//创建一个存放bytes字节的缓冲区
	dataBuff := bytes.NewBuffer(binaryData)

	//直接压head信息，得到dataLen和msgId
	msg := &Message{}

	//1.读取dataLen
	if err := binary.Read(dataBuff,binary.LittleEndian,&msg.DataLen); err != nil{
		return nil,err
	}
	//2.读MsgID
	if err := binary.Read(dataBuff,binary.LittleEndian,&msg.Id); err != nil{
		return nil,err
	}
	//判断dataLen 是否已经超出了我们允许的最大的包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize  {
		return nil,errors.New("too Large msg data recv!")
	}
	//这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据
	return msg,nil
}





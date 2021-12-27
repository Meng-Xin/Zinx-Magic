package znet

import (
	"fmt"
	"strconv"
	"zinx-lwh/ziface"
)

/*
	消息处理模块的实现
 */

type MsgHandle struct {
	//接口集合，存放每个MsgID对应的处理方法
	Apis map[uint32] ziface.IRouter

}
// NewMsgHandle 初始化/创建MsgHandle 方法
func NewMsgHandle() *MsgHandle  {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

// DoMsgHandler 调度/执行对应的Router消息处理方法
func (m *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	//1.从Request中找到msgID
	handler,ok :=m.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID=",request.GetMsgID(),"is NOT Found！You Need Register")
	}
	//调用三个 方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}
// AddRouter 为消息添加具体的处理逻辑
func (m *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	//1.判断 当前msg绑定的API处理方法是否已经存在
	if _,ok := m.Apis[msgID]; ok {
		//当前 id已经被注册了
		panic("repeat api,msgID="+strconv.Itoa(int(msgID)))
	}
	//2.添加msg与api的绑定关系
	m.Apis[msgID] = router
	fmt.Println("Add api MsgID=",msgID,"succ!")
}



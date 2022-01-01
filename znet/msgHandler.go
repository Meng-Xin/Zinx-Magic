package znet

import (
	"fmt"
	"strconv"
	"zinx-lwh/utils"
	"zinx-lwh/ziface"
)

/*
	消息处理模块的实现
*/

type MsgHandle struct {
	//接口集合，存放每个MsgID对应的处理方法
	Apis map[uint32]ziface.IRouter
	//负责Worker 取任务的消息队列
	TaskQueue []chan ziface.IRequest
	//业务工作Worker池的worker数量
	WorkerPoolSize uint32
}

// NewMsgHandle 初始化/创建MsgHandle 方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, //从全局配置中读取用户设置的Goroutine数量
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// DoMsgHandler 调度/执行对应的Router消息处理方法
func (m *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	//1.从Request中找到msgID
	handler, ok := m.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID=", request.GetMsgID(), "is NOT Found！You Need Register")
		return
	}
	//调用三个 方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter 为消息添加具体的处理逻辑
func (m *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	//1.判断 当前msg绑定的API处理方法是否已经存在
	if _, ok := m.Apis[msgID]; ok {
		//当前 id已经被注册了
		panic("repeat api,msgID=" + strconv.Itoa(int(msgID)))
	}
	//2.添加msg与api的绑定关系
	m.Apis[msgID] = router
	fmt.Println("Add api MsgID=", msgID, "succ!")
}

// StartWorkerPool 启动一个Worker工作池 (开启工作池动作只能发生一次，一个Zinx框架只能拥有一个工作池)
func (m *MsgHandle) StartWorkerPool() {
	//根据workerPoolSize 分别开启Worker ，每一个Worker用一个go来承载
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		//一个worker被启动
		//1.给当前的worker对用的channel 消息队列开辟空间 【有缓冲chan 能够缓冲1024个任务】
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//2.启动当前的Worker ，阻塞等待消息从channel传递进来
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

// StartOneWorker 启动一个Worker工作流程
func (m *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, "is start...")

	//不断的阻塞等待对应消息队列的消息
	for {
		select {
		//如果有消息过来，出列的就是一个客户端的Request,执行当前Request所绑定的业务
		case request := <-taskQueue:
			m.DoMsgHandler(request)
		}
	}
}

// SendMsgToTaskQueue 将消息交给TaskQueue, 由Worker进行处理
func (m *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//1.将消息平均分配给Worker
	//根据客户端建立的ConnID来进行分配
	workerID := request.GetConnection().GetConnID() % m.WorkerPoolSize
	fmt.Println("Add ConnID =", request.GetConnection().GetConnID(), "request MsgID", request.GetMsgID(), "to WorkerID", workerID)
	//2.将消息发送给对应的worker的TaskQueue即可
	m.TaskQueue[workerID] <- request
}

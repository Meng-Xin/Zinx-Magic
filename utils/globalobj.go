package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx-lwh/ziface"
)

/*
	存储一切有关Zinx框架的全局参数，供其他模块使用
	一些参数是可以通过zinx.json由永固进行配置的
*/

type GlobalObj struct {
	//Server
	TcpServer ziface.IServer //当前Zinx全局的Server对象
	Host      string         //当前服务器主机监听的IP
	TcpPort   int            //当前服务器足迹监听的端口号
	Name      string         //当前服务器的名称

	//Zinx
	Version          string //当前Zinx的版本号
	MaxConn          int    //当前服务器主机允许的最大连接数
	MaxPackageSize   uint32 //当前Zinx框架数据包的最大值
	WorkerPoolSize   uint32 // Worker工作池的消息队列个数
	MaxWorkerTaskLen uint32 //每个Worker对应的消息队列的任务的最大数量
}

// GlobalObject 定义一个全局GloablObj
var GlobalObject *GlobalObj

//Reload 从 zinx.json 中加载用于自定义的参数
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	//将json文件数据解析到结构体中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}

}

// 提供一个init方法，初始化GlobalObject
func init() {
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "V0.10",
		Host:             "0.0.0.0",
		TcpPort:          8999,
		MaxConn:          3,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}

	//应该尝试从conf/zinx.json中加载一些用户自定义参数
	GlobalObject.Reload()
}

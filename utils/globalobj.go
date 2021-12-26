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
	Version        string //当前Zinx的版本号
	MaxConn        int    //当前服务器主机允许的最大连接数
	MaxPackageSize uint32 //当前Zinx框架数据包的最大值
}
// 定义一个全局GloablObj
var GlobalObject *GlobalObj

//从 zinx.json 中加载用于自定义的参数
func (g *GlobalObj)Reload()  {
	data,err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	//将json文件数据解析到结构体中
	err = json.Unmarshal(data,&GlobalObject)
	if err != nil {
		panic(err)
	}

}
// 提供一个init方法，初始化GlobalObject
func init()  {
	GlobalObject = &GlobalObj{
		Name: "ZinxServerApp",
		Version: "V0.4",
		Host: "0.0.0.0",
		TcpPort: 8999,
		MaxConn: 1000,
		MaxPackageSize: 4096,
	}

	//应该尝试从conf/zinx.json中加载一些用户自定义参数
	//GlobalObject.Reload()
}
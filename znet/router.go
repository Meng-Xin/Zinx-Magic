package znet

import "zinx-lwh/ziface"

//实现router时，线迁出这个BaseRouter基类，让后根据需要对这个基类的方法进行重写就好了。
type BaseRouter struct {}
/*
 1.这里之所以BaseRouter的方法都为空是因为有的Router
不希望有PreHandle、PostHandle这两个业务，所以Router
全部继承BaseRouter的好处就是，不需要实现PreHandle、PostHandle
*/

// PreHandle 处理conn 业务之前的钩子方法。
func (b BaseRouter) PreHandle(request ziface.IRequest) {}
// Handle 处理conn 业务的主方法Hook
func (b BaseRouter) Handle(request ziface.IRequest) {}
// PostHandle 处理conn 业务的之后的方法Hook
func (b BaseRouter) PostHandle(request ziface.IRequest) {}




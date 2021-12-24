package ziface

// IRequest 接口：实际上是把客户端请求的链接信息，和 请求数据 包装到一个Request中，
type IRequest interface {
	// GetConnection 得到链接
	GetConnection() IConnection
	// GetData 客户端请求数据
	GetData()	[]byte
}

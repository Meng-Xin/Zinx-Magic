package ziface


type IDataPack interface {
	// GetHeadLen 获取爆头的长度方法
	GetHeadLen() uint32
	// Pack 封包方法
	Pack(message IMessage) ([]byte,error)
	// Unpack 拆包方法
	Unpack([]byte)(IMessage,error)
}

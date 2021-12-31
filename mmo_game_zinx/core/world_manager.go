package core

import "sync"

type WorldManager struct {
	//AOIManager 当前世界地图AOI 的管理模块
	AoiMgr *AOIManager
	//当前全部在线的Players集合
	Players map[int32] *Player
	//保护PLayers集合的锁
	pLock sync.RWMutex
}

// WorldMgrObj 提供一个对外开放的世界管理模块的句柄
var WorldMgrObj *WorldManager
//初始化方法
func init()  {
	WorldMgrObj = &WorldManager{
		//创建世界AOI地图规划
		AoiMgr: NewAOIManager(AOI_MIN_X,AOI_MAX_X,AOI_CNTS_X,AOI_MIN_Y,AOI_MAX_Y,AOI_CNTS_Y),
		//初始化Player集合
		Players :make(map[int32]*Player),
	}
}
// AddPlayer 添加一个 玩家
func (wm *WorldManager)AddPlayer(player *Player)  {
	wm.pLock.Lock()
	wm.Players[player.Pid] = player
	wm.pLock.Unlock()

	//将player添加到 , AOIManager中
	wm.AoiMgr.AddToGridByPos(int(player.Pid),player.X,player.Y)
}
// RemovePlayerByPid 删除一个玩家
func (wm *WorldManager)RemovePlayerByPid(pid int32)  {
	//得到当前玩家
	player := wm.Players[pid]
	//将玩家从AOIManager中删除
	wm.AoiMgr.RemoveFromGridByPos(int(pid),player.X,player.Y)

	//将玩家从世界管理中删除
	wm.pLock.Lock()
	delete(wm.Players,pid)
	wm.pLock.Unlock()
}
// GetPlayerByPid 通过玩家ID 查询Player对象
func (wm *WorldManager) GetPlayerByPid (pid int32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	return wm.Players[pid]
}
// GetAllPlayers 获取全部的在线玩家
func (wm *WorldManager) GetAllPlayers (pid int32) []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	players := make([]*Player,0)
	//添加到切片中
	for _,p := range wm.Players{
		players = append(players, p)
	}
	return players
}
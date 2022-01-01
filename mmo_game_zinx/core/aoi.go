package core

import "fmt"

//定义一些AOI的边界值
const (
	AOI_MIN_X  int = 85
	AOI_MAX_X  int = 410
	AOI_CNTS_X int = 10
	AOI_MIN_Y  int = 75
	AOI_MAX_Y  int = 400
	AOI_CNTS_Y int = 20
)

/*
	AOI 区域管理模块
*/

type AOIManager struct {
	//区域的左边界坐标
	MinX int
	//区域的右边界坐标
	MaxX int
	//X方向格子的数量
	CntX int
	//区域的上边界坐标
	MinY int
	//区域的下边界坐标
	MaxY int
	//Y方向格子的数量
	CntY int
	//当前区域中有那些格子map-key=格子的ID ,value= 格子的对象
	grids map[int]*Grid
}

/*
	初始化一个AOI 区域管理模块
*/

func NewAOIManager(minX, maxX, cntX, minY, maxY, cntY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		MinY:  minY,
		MaxY:  maxY,
		CntX:  cntX,
		CntY:  cntY,
		grids: make(map[int]*Grid),
	}
	//给AOI 初始化区域的格子所有的格子进行编号和初始化
	for y := 0; y < cntY; y++ {
		for x := 0; x < cntX; x++ {
			//计算 格子id ，根据x,y编号
			//格子编号: id = idy * cntX + idx
			gid := y*cntX + x
			//初始化gid 格子
			aoiMgr.grids[gid] = NewGrid(gid,
				aoiMgr.MinX+x*aoiMgr.gridWidth(),
				aoiMgr.MinX+(x+1)*aoiMgr.gridWidth(),
				aoiMgr.MinY+y*aoiMgr.gridLength(),
				aoiMgr.MinY+(y+1)*aoiMgr.gridLength(),
			)
		}
	}
	return aoiMgr
}

// gridWidth 得到每个格子在X轴方向的宽度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntX
}

// gridLength 得到每个格子在Y轴方向的高度
func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.CntY
}

// String 打印格子信息
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager:\n MinX:%d,MaxX:%d,MinY:%d,MaxY:%d,CntX:%d,CntY:%d",
		m.MinX, m.MaxX, m.MinY, m.MaxY, m.CntX, m.CntY,
	)
	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}
	return s
}

// GetSurroundGridsByGid 根据格子GID得到周边九宫格格子集合
func (m *AOIManager) GetSurroundGridsByGid(gID int) (grids []*Grid) {
	//1.判断GID 是否在AOIManager中
	if _, ok := m.grids[gID]; !ok {
		return
	}
	//初始话grids返回值 切片
	grids = append(grids, m.grids[gID])
	//需要gID得到左边是否有格子? 右边是否有格子
	idx := gID % m.CntX
	//需要通过gID 得到当前格子X轴的编号-- idx = id % cntX
	//判断idx编号左边是否有格子，如果有 放在gidX 集合中
	if idx > 0 {
		grids = append(grids, m.grids[gID-1])
	}
	//判断idx编号右边是否有格子，如果有 放在gidX 集合中
	if idx < m.CntX-1 {
		grids = append(grids, m.grids[gID+1])
	}
	//遍历gidX 集合中每个格子的gid
	gidsX := make([]int, 0, len(grids))
	for _, v := range grids {
		gidsX = append(gidsX, v.GID)
	}
	//遍历gidsX 集合中每个格子的gid
	for _, v := range gidsX {
		idy := v / m.CntY
		//gid 上边是否还有格子
		if idy > 0 {
			grids = append(grids, m.grids[v-m.CntX])
		}
		//gid 下边是否还有格子
		if idy < m.CntY-1 {
			grids = append(grids, m.grids[v+m.CntX])
		}
	}
	return
}

// GetGidByPos 通过x,y横纵轴坐标得到当前的GID 格子编号
func (m *AOIManager) GetGidByPos(x, y float32) int {
	idx := (int(x) - m.MinX) / m.gridWidth()
	idy := (int(y) - m.MinY) / m.gridLength()

	return idy*m.CntX + idx
}

// GetPidByPos 通过横纵坐标得到周边九宫格内全部的PlayerIDS
func (m *AOIManager) GetPidByPos(x, y float32) (playerIDs []int) {
	//得到当前玩家的GID 格子id
	gID := m.GetGidByPos(x, y)
	//通过Gid 得到周边九宫格信息
	grids := m.GetSurroundGridsByGid(gID)
	//将九宫格的信息李的全部Player的id 累加到 playerIDs
	for _, v := range grids {
		playerIDs = append(playerIDs, v.GetPlayerIDs()...)
		fmt.Printf("===>grid ID :%d,pid :%v", v.GID, v.GetPlayerIDs())
	}
	return
}

// AddPidToGrid 添加一个PlayerID到一个格子中
func (m *AOIManager) AddPidToGrid(pID, gID int) {
	m.grids[gID].Add(pID)
}

// RemovePidFromGrid 移除GID 获取全部的PlayerID
func (m *AOIManager) RemovePidFromGrid(pID, gId int) {
	m.grids[gId].Remove(pID)
}

// GetPidByGid 通过GID 获取全部的PlayerID
func (m *AOIManager) GetPidByGid(gID int) (playerIDs []int) {
	playerIDs = m.grids[gID].GetPlayerIDs()
	return
}

// AddToGridByPos 通过坐标将Player添加到一个格子中
func (m *AOIManager) AddToGridByPos(pID int, x, y float32)  {
	fmt.Println("能进来吗？")
	gID := m.GetGidByPos(x,y)
	grID := m.grids[gID]
	grID.Add(pID)
	fmt.Println("能否结束？")
}

//
// RemoveFromGridByPos 通过坐标把一个Player从一个格子中删除
func (m *AOIManager) RemoveFromGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	grid := m.grids[gID]
	grid.Remove(pID)

}

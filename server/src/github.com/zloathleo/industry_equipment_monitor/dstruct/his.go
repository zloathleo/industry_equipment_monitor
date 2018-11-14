package dstruct

//历史曲线切片
type PointsValueSnap map[string]*HDas

//历史曲线数据map
type PointsValueHistoryMap map[string][]*HDas

//历史曲线返回值
type PointsValueHistory struct {
	XAxis  []int64               `json:"xAxis"`
	Series PointsValueHistoryMap `json:"series"`
}

//初始化历史返回值
func NewPointsValueHistory(pointsArray []string, hisReqParam *HisReqParam) *PointsValueHistory {
	hisCount := hisReqParam.Count
	interval := hisReqParam.Interval

	history := &PointsValueHistory{}
	//先初始化数据结果
	//时间戳列表
	xAxis := make([]int64, hisCount, hisCount)
	for i := 0; i < hisCount; i++ {
		xAxis[i] = hisReqParam.From.Unix() + int64(i*interval)
	}
	history.XAxis = xAxis

	//历史数据结果集
	history.Series = make(PointsValueHistoryMap)
	for _, point := range pointsArray {
		history.Series[point] = make([]*HDas, hisCount, hisCount)
	}
	return history
}

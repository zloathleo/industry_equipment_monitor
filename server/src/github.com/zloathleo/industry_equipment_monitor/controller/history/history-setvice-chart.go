package history

import (
	"bytes"
	"github.com/zloathleo/industry_equipment_monitor/common/logger"
	. "github.com/zloathleo/industry_equipment_monitor/dstruct"
	"github.com/zloathleo/industry_equipment_monitor/his"
	"github.com/zloathleo/industry_equipment_monitor/utils"
	"strconv"
	"time"
)

/**
格式化历史对齐数据
dur 为枚举值[600...86400]
间隔interval 为枚举值[0,2...60]
300 	5分钟   interval=1     300点+1  即获取所有
600 	10分钟   interval=0     600点+1  即获取所有
1800	30分钟   interval=2     900点+1
3600 	60分钟   interval=5     720点+1
14400	4小时	 interval=15    960点+1
28800 	8小时	 interval=30    960点+1
86400 	24小时   interval=60    1440点+1
604800 	168小时  interval=600    1008点+1
*/
func fetchHistoryChartData(pointsArray []string, to int64, dur int, interval int) *PointsValueHistory {
	hisReqParam := his.NewHisReqParam(to, dur, interval)

	if hisReqParam == nil {
		logger.Warnf("request his param is err to[%d], dur[%d]", to, dur)
		return nil
	}

	times := hisReqParam.Times
	if times == nil || len(times) == 0 {
		logger.Warnf("request his param is err to[%d], dur[%d]", to, dur)
		return nil
	}
	pointsValueHistoryMap := NewPointsValueHistory(pointsArray, hisReqParam)

	//历史请求总的上下文
	hisReqContext := &HisReqContext{Count: hisReqParam.Count, Index: 0}

	//按天分文件查询
	//times[0].From = times[0].From - 60
	//fillOneDayHistoryChartData(pointsArray, xAxis, times[0], hisMap, hisReqContext)
	//
	//for i:=1;i<len(times);i++{
	//	item := times[i]
	//	fillOneDayHistoryChartData(pointsArray, xAxis, item, hisMap, hisReqContext)
	//}
	for _, item := range times {
		begin := time.Now().UnixNano()
		fillOneDayHistoryChartData(pointsArray, item, pointsValueHistoryMap, hisReqContext)
		logger.Debugf("fillOneDayHistoryChartData const: %d ms. %v", (time.Now().UnixNano()-begin)/int64(1000000), item)
	}
	return pointsValueHistoryMap

}

//只能查看from当天
func fillOneDayHistoryChartData(pointsArray []string, paramTime *HisReqParamTime, historyMap *PointsValueHistory, hisReqContext *HisReqContext) {

	fromTime := time.Unix(paramTime.From, 0)
	toTime := time.Unix(paramTime.To, 0)

	fromHour := utils.GetTimeIntHour(fromTime).Hour()
	toHour := utils.GetTimeIntHour(toTime).Hour()

	blockForm := fromTime
	blockTo := toTime
	//hisCount := hisReqParam.Count
	interval := paramTime.Interval

	hourIndex := fromHour
	//从第一个小时开始
	for ; hourIndex < toHour; hourIndex++ {
		blockTo = utils.GetTimeEndOfHour(blockForm)

		//当前table 需要查找的数量
		tableDataCount := int(utils.GetTimeNextIntHour(blockForm).Unix()-blockForm.Unix()) / interval

		reqParamTime := NewHisReqParamTime(blockForm.Unix(), blockTo.Unix(), tableDataCount, interval)
		his.SelectSingleTableHistoryData(pointsArray, reqParamTime, historyMap, hisReqContext)

		hisReqContext.Index = hisReqContext.Index + tableDataCount
		blockForm = utils.GetTimeNextIntHour(blockTo)
	}

	tableDataCount := int(paramTime.To-blockForm.Unix())/interval + 1
	his.SelectSingleTableHistoryData(pointsArray, NewHisReqParamTime(blockForm.Unix(), paramTime.To, tableDataCount, interval), historyMap, hisReqContext)
	hisReqContext.Index = hisReqContext.Index + tableDataCount

}

//输出历史曲线Json数据
func renderChartHistoryJson(historyMap *PointsValueHistory) *bytes.Buffer {
	var builder bytes.Buffer

	//root
	builder.WriteString("{ ")

	//xAxis
	builder.WriteString("\"xAxis\":[")
	for _, ts := range historyMap.XAxis {
		//builder.WriteString(strconv.FormatInt(ts,10) + ",")
		builder.WriteString("\"" + utils.GetIntTimeString(ts) + "\",")
	}
	builder.Truncate(builder.Len() - 1)
	builder.WriteString("],")

	//his data
	builder.WriteString("\"series\":{ ")
	for pn, dasArray := range historyMap.Series {
		builder.WriteString("\"" + pn + "\":[ ")

		for _, das := range dasArray {
			if das != nil {
				builder.WriteString(strconv.FormatFloat(das.V, 'f', 2, 64) + ",")
			} else {
				builder.WriteString("null,")
			}
		}
		builder.Truncate(builder.Len() - 1)

		builder.WriteString(" ],")
	}
	builder.Truncate(builder.Len() - 1)
	builder.WriteString(" }")

	//root
	builder.WriteString(" }")
	return &builder
}

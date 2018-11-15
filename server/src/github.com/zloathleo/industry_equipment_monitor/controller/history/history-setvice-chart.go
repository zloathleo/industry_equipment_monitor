package history

import (
	"time"
	"bytes"
	"strconv"
	"github.com/zloathleo/industry_equipment_monitor/common/logger"
	"github.com/zloathleo/industry_equipment_monitor/utils"
	. "github.com/zloathleo/industry_equipment_monitor/dstruct"
	"github.com/zloathleo/industry_equipment_monitor/his"
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
func fetchHistoryChartData(pointsArray []string, to int64, dur int, interval int) (map[string][]*HDas, []int64) {
	if pointsArray == nil || len(pointsArray) == 0 {
		return nil, nil
	}
	hisReqParam := his.FormatChartHisReqParam(to, dur, interval)

	if hisReqParam == nil {
		logger.Warnf("request his param is err to[%d], dur[%d]", to, dur)
		return nil, nil
	}

	times := hisReqParam.Times
	if times == nil || len(times) == 0 {
		logger.Warnf("request his param is err to[%d], dur[%d]", to, dur)
		return nil, nil
	} else {
		hisCount := hisReqParam.Count
		interval := hisReqParam.Interval
		//先初始化数据结果
		//时间戳列表
		xAxis := make([]int64, hisCount, hisCount)
		for i := 0; i < hisCount; i++ {
			xAxis[i] = hisReqParam.From.Unix() + int64(i*interval)
		}
		//历史数据结果集
		var hisMap = make(map[string][]*HDas)
		for _, point := range pointsArray {
			hisMap[point] = make([]*HDas, hisCount, hisCount)
		}
		//历史请求总的上下文
		hisReqContext := &HisReqContext{Count: hisCount, Index: 0}

		//按天分文件查询
		for _, item := range times {
			begin := time.Now().UnixNano()
			fillOneDayHistoryChartData(pointsArray, xAxis, item, hisMap, hisReqContext)
			logger.Debugf("fillOneDayHistoryChartData const: %d ms. %v", (time.Now().UnixNano()-begin)/int64(1000000), item)
		}
		return hisMap, xAxis
	}

}

//只能查看from当天
func fillOneDayHistoryChartData(pointsArray []string, xAxis []int64, paramTime *HisReqParamTime, hisMap map[string][]*HDas, hisReqContext *HisReqContext) {

	fromTime := time.Unix(paramTime.From, 0)
	toTime := time.Unix(paramTime.To, 0)

	fromHour := utils.GetTimeIntHour(fromTime).Hour()
	toHour := utils.GetTimeIntHour(toTime).Hour()

	blockForm := fromTime
	blockTo := toTime
	//hisCount := hisReqParam.Count
	interval := paramTime.Interval

	i := fromHour
	//从第一个小时开始
	for ; i < toHour; i++ {
		blockTo = utils.GetTimeEndOfHour(blockForm)

		//当前table 需要查找的数量
		tableDataCount := int(utils.GetTimeNextIntHour(blockForm).Unix()-blockForm.Unix()) / interval
		begin := time.Now().UnixNano()
		his.SelectSingleTableHistoryData(pointsArray, blockForm, i, blockForm.Unix(), blockTo.Unix(), interval, tableDataCount, hisMap, hisReqContext)
		logger.Debugf("SelectSingleTableHistoryData const: %d ms", (time.Now().UnixNano()-begin)/int64(1000000))

		hisReqContext.Index = hisReqContext.Index + tableDataCount
		logger.Debugf("请求 小时分割  tableDataCount %v || %v || %d", blockForm, blockTo, tableDataCount)
		blockForm = utils.GetTimeNextIntHour(blockTo)
	}

	tableDataCount := int(paramTime.To-blockForm.Unix())/interval + 1
	logger.Debugf("req time tableDataCount %v || %v || %d", blockForm, paramTime.To, tableDataCount)
	his.SelectSingleTableHistoryData(pointsArray, blockForm, toHour, blockForm.Unix(), paramTime.To, paramTime.Interval, tableDataCount, hisMap, hisReqContext)
	hisReqContext.Index = hisReqContext.Index + tableDataCount

}

//输出历史曲线Json数据
func renderChartHistoryJson(hisMap map[string][]*HDas, xAxis []int64) *bytes.Buffer {
	var builder bytes.Buffer

	//root
	builder.WriteString("{ ")

	//xAxis
	builder.WriteString("\"xAxis\":[")
	for _, ts := range xAxis {
		//builder.WriteString(strconv.FormatInt(ts,10) + ",")
		builder.WriteString("\"" + utils.GetIntTimeString(ts) + "\",")
	}
	builder.Truncate(builder.Len() - 1)
	builder.WriteString("],")

	//his data
	builder.WriteString("\"series\":{ ")
	for pn, dasArray := range hisMap {
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

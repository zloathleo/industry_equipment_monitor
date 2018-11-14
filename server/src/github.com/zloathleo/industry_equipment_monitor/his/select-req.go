package his

import (
	"github.com/zloathleo/industry_equipment_monitor/common/logger"
	"github.com/zloathleo/industry_equipment_monitor/dstruct"
	"github.com/zloathleo/industry_equipment_monitor/utils"
	"time"
)

/**
格式化当天历史请求
只能查看from当天
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
func NewHisReqParam(toInt int64, durInt int, interval int) *dstruct.HisReqParam {
	//起始时间
	var fromInt int64

	//时长
	switch durInt {
	case 300:
		{
			//不改
			fromInt = toInt - 300
			if interval == 0 {
				interval = 1
			}
			break
		}
	case 600:
		{
			//不改
			fromInt = toInt - 600
			if interval == 0 {
				interval = 1
			}
			break
		}
	case 1800:
		{
			//30分钟
			fromInt = toInt - 1800
			if interval == 0 {
				interval = 2
			}
			break
		}
	case 3600:
		{
			//60分钟
			fromInt = toInt - 3600
			if interval == 0 {
				interval = 5
			}

			//对齐时间
			fromInt = utils.GetTimeInt5Sec(time.Unix(fromInt, 0)).Unix()
			toInt = fromInt + 3600

			break
		}
	case 14400:
		{
			//4小时
			fromInt = toInt - 14400
			if interval == 0 {
				interval = 15
			}
			//对齐时间
			fromInt = utils.GetTimeInt15Sec(time.Unix(fromInt, 0)).Unix()
			toInt = fromInt + 14400

			break
		}
	case 28800:
		{
			//8小时
			fromInt = toInt - 28800
			if interval == 0 {
				interval = 30
			}

			//对齐时间
			fromInt = utils.GetTimeInt30Sec(time.Unix(fromInt, 0)).Unix()
			toInt = fromInt + 28800
			break
		}
	case 86400:
		{
			//24小时
			fromInt = toInt - 86400
			if interval == 0 {
				interval = 60
			}

			//对齐时间
			fromInt = utils.GetTimeIntMin(time.Unix(fromInt, 0)).Unix()
			toInt = fromInt + 86400
			break
		}
	case 604800:
		{
			//24小时*7
			fromInt = toInt - 604800
			if interval == 0 {
				interval = 600
			}

			//对齐时间
			fromInt = utils.GetTimeInt10Min(time.Unix(fromInt, 0)).Unix()
			toInt = fromInt + 604800
			break
		}
	default:
		{
			return nil
		}
	}

	fromTime := time.Unix(fromInt, 0)
	toTime := time.Unix(toInt, 0)

	//按天分割的时间片段
	var times []*dstruct.HisReqParamTime
	//需要分文件查询
	fromStepTimeInt := fromInt //int64
	toStepTime := utils.GetTimeDayEnd(fromTime)

	for !toStepTime.After(toTime) {
		itemCount := int(toStepTime.Unix()-fromInt) / interval
		toStepTimeInt := fromInt + int64(interval*itemCount)

		//第二天 的下一个 时间戳
		nextFromStepTimeInt := toStepTimeInt + int64(interval)

		//add
		hisReqParamTime := &dstruct.HisReqParamTime{From: fromStepTimeInt, To: toStepTimeInt, Count: itemCount, Interval: interval}
		times = append(times, hisReqParamTime)

		fromStepTimeInt = nextFromStepTimeInt
		toStepTime = utils.GetTimeDayEnd(time.Unix(fromStepTimeInt, 0))

	}

	hisReqParamTime := &dstruct.HisReqParamTime{From: fromStepTimeInt, To: toInt, Count: int(toInt-fromStepTimeInt)/interval + 1, Interval: interval}
	logger.Debug("请求 日期分割  时间 step last : ", utils.GetIntTimeString(fromStepTimeInt), utils.GetIntTimeString(toInt))
	times = append(times, hisReqParamTime)

	return &dstruct.HisReqParam{From: fromTime, To: toTime, Times: times, Dur: durInt, Interval: interval, Count: durInt/interval + 1}
}

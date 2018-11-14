package dstruct

import "time"

type HisReqParam struct {
	From     time.Time
	To       time.Time
	Times    []*HisReqParamTime
	Dur      int // 时长
	Interval int // 时间戳间隔[1...60]
	Count    int // Dur/Interval + 1
}

type HisReqParamTime struct {
	From     int64
	To       int64
	Count    int // Dur/Interval + 1
	Interval int // 时间戳间隔[1...60]
}

func NewHisReqParamTime(f int64, t int64, c int, i int) *HisReqParamTime {
	return &HisReqParamTime{From: f, To: t, Count: c, Interval: i}
}

//历史请求总的上下文
type HisReqContext struct {
	Count int //历史时间戳总数
	Index int //历史时间戳序号
}

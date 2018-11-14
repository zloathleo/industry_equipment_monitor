package dstruct

type Alarm struct {
	Module string `json:"module"`
	Err string `json:"err"`
	TimeStamp int64  `json:"timestamp"`
}

type PushDas struct {
	Rows      []*Das `json:"rows"`
	TimeStamp uint64 `json:"timestamp"`
}

type Das struct {
	PointName string  `json:"name,omitempty"`
	Value     float64 `json:"value"`
	TimeStamp uint64  `json:"timestamp"`
}

type HDas struct {
	V float64 `json:"v"`
	T int64  `json:"t"`
}





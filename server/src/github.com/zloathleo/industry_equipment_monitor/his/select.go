package his

import (
	"database/sql"
	"github.com/zloathleo/industry_equipment_monitor/common/logger"
	"github.com/zloathleo/industry_equipment_monitor/db"
	. "github.com/zloathleo/industry_equipment_monitor/dstruct"
	"github.com/zloathleo/industry_equipment_monitor/systemalarm"
	"github.com/zloathleo/industry_equipment_monitor/utils"
	"strconv"
	"strings"
	"time"
)

func getSelectPointsString(points []string) string {
	a := make([]string, len(points))
	for i := 0; i < len(points); i++ {
		a[i] = "'" + points[i] + "'"
	}
	return strings.Join(a, ",")
}

/**
查找指定文件 指定表格的历史
*/
func SelectSingleTableHistoryData(points []string, reqParamTime *HisReqParamTime, historyMap *PointsValueHistory, hisReqContext *HisReqContext) {
	begin := time.Unix(reqParamTime.From, 0)

	currentConnect := db.GetAnyConnect(begin)
	if currentConnect == nil {
		//历史不存在
		logger.Infof("can find connect %s.", utils.GetDayString(begin))
		return
	}
	tableName := findHisTable(begin)
	pointsString := getSelectPointsString(points)

	selectSql := "select * from " + tableName + " t where t.pn in (" + pointsString + ") and t.t >= " + strconv.FormatInt(begin.Unix(), 10) + " and t.t <= " + strconv.FormatInt(reqParamTime.To, 10) + " ORDER BY t.t"
	logger.Debug(selectSql)

	var readLock = todayDbLock
	var readDb = todayDb
	//如果不是查询今天
	if currentConnect.Name != todayDbName {
		readLock = currentConnect.DBLock
		readDb = currentConnect.DB
	}
	readLock.RLock()
	rows, err := readDb.Query(selectSql, 1)

	defer func() {
		readLock.RUnlock()
		if rows != nil {
			rows.Close()
		}
	}()

	if err == nil && rows != nil {
		foreachRow(rows, points, reqParamTime, historyMap, hisReqContext)
	} else {
		systemalarm.AddSystemAlarm(err)
		return
	}

}

func foreachRow(rows *sql.Rows, points []string, reqParamTime *HisReqParamTime, historyMap *PointsValueHistory, hisReqContext *HisReqContext) {

	var pn string
	var t int64 //查询出的值时刻点
	var value float64

	//上一次的Index
	pointJumpMap := make(map[string]int)
	//上一次的值
	pointValueMap := make(map[string]*HDas)
	for _, pn := range points {
		//pointBeginMap[pn] = uint64(from)
		pointJumpMap[pn] = 0
		pointValueMap[pn] = nil
	}

	for rows.Next() {
		err := rows.Scan(&pn, &t, &value)
		if err != nil {
			systemalarm.AddSystemAlarm(err)
			return
		} else {
			handleRecode(pn, t, value, reqParamTime, historyMap, hisReqContext, pointJumpMap, pointValueMap)
		}
	}

	for pointName, lastDas := range pointValueMap {
		if lastDas != nil {
			//最后一个数值补齐  补齐后续1分钟
			lastEndTime := lastDas.T + 60
			if lastEndTime > reqParamTime.To {
				lastEndTime = reqParamTime.To
			}
			jump := int(lastEndTime-reqParamTime.From) / reqParamTime.Interval
			lastJump := pointJumpMap[pointName]

			//从上次Index开始补齐值
			for j := lastJump + 1; j <= jump; j++ {
				jIdx := hisReqContext.Index + j
				historyMap.Series[pointName][jIdx] = &HDas{V:lastDas.V}
			}
		}
	}

}

/**
pn 点名
t 点时间戳
value 值
fromUInt64 开始
interval 间隔
tableDataCount
hisMap历史结果集
hisReqContext历史请求上下文
*/
func handleRecode(pn string, t int64, value float64, reqParamTime *HisReqParamTime, historyMap *PointsValueHistory, hisReqContext *HisReqContext, pointJumpMap map[string]int, pointValueMap map[string]*HDas) {

	remainder := int(t-reqParamTime.From) % reqParamTime.Interval
	jump := int(t-reqParamTime.From) / reqParamTime.Interval
	newIndex := hisReqContext.Index + jump

	lastJump := pointJumpMap[pn]
	lastValue := pointValueMap[pn]
	if lastValue == nil {

	} else {
		//从上次Index开始补齐值
		maxAppend := jump
		if t-reqParamTime.From > 60 {
			maxAppend = lastJump + 60 / reqParamTime.Interval
		}

		for j := lastJump + 1; j <= maxAppend; j++ {
			jIdx := hisReqContext.Index + j
			historyMap.Series[pn][jIdx] = &HDas{V: lastValue.V}
		}
	}

	//余数==0代表就在需要的时间戳伤
	if remainder == 0 {
		historyMap.Series[pn][newIndex] = &HDas{V: value}
	} else {
	}

	pointJumpMap[pn] = jump
	pointValueMap[pn] = &HDas{V: value, T: t}
}

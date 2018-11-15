package his

import (
	"fmt"
	"github.com/zloathleo/industry_equipment_monitor/common/logger"
	"strings"
	"time"
)

func saveRealtime(cacheCopy map[string]float64, ignore bool) {
	count := len(cacheCopy)
	if count == 0 {
		return
	}
	now := time.Now()
	timeStamp := now.Unix()
	tableName := findHisTable(now)

	valueStrings := make([]string, 0, count)
	valueArgs := make([]interface{}, 0, count*3)
	for pn, value := range cacheCopy {
		valueStrings = append(valueStrings, "(?, ?, ?)")
		valueArgs = append(valueArgs, pn)
		valueArgs = append(valueArgs, timeStamp)
		valueArgs = append(valueArgs, value)
	}
	sqlInsertCmd := "INSERT"
	if ignore {
		sqlInsertCmd = "INSERT OR IGNORE"
	}
	stmt := fmt.Sprintf(sqlInsertCmd+" INTO "+tableName+" (pn, t, value) VALUES %s", strings.Join(valueStrings, ","))

	logger.Debug(stmt)
	//执行插入
	todayDbLock.Lock()
	_, err := todayDb.Exec(stmt, valueArgs...)
	todayDbLock.Unlock()

	if err != nil {
		logger.Warnf("insert real time data err [%s].", err.Error())
	}

}

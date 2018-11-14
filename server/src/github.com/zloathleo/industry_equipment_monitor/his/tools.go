package his

import (
	"time"
	"fmt"
)

//查找合适的历史表格
func findHisTable(timeStamp time.Time) string {
	//1539680150
	return fmt.Sprintf("his24_%d", timeStamp.Hour())
}

//查找合适的历史表格
func findHisTableByIndex(index int) string {
	return fmt.Sprintf("his24_%d", index)
}

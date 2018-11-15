package his

import (
	"github.com/zloathleo/industry_equipment_monitor/common/logger"
)

func InitHis(){
	initDB()
	initCorn()
	logger.Warnln("his init ok.")
}

package common

import (
	"github.com/zloathleo/industry_equipment_monitor/common/config"
	"github.com/zloathleo/industry_equipment_monitor/common/logger"
)

func InitCommon() {
	logger.Init()
	config.Init()
	logger.SetLevel(config.AppConfig.App.LogLevel)
	logger.Warnln("common init ok.")
}  
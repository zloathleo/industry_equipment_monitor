package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/zloathleo/industry_equipment_monitor/controller/scheme"
	"github.com/zloathleo/industry_equipment_monitor/controller/realtime"
	"github.com/zloathleo/industry_equipment_monitor/controller/history"
)

func InitController(dasGroup *gin.RouterGroup) {
	scheme.InitSchemeController(dasGroup)
	realtime.InitRealtimeController(dasGroup)
	history.InitHistoryController(dasGroup)
}
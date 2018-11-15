package main

import (
	"github.com/zloathleo/industry_equipment_monitor/common/logger"

	"github.com/zloathleo/industry_equipment_monitor/common"
	"github.com/zloathleo/industry_equipment_monitor/httpr"

	"os"
	"os/signal"

	//"github.com/zloathleo/industry_equipment_monitor/history2"
	"github.com/zloathleo/industry_equipment_monitor/his"
	"github.com/zloathleo/industry_equipment_monitor/memcache"
	"github.com/zloathleo/industry_equipment_monitor/pointscheme"
)

func main() {
	common.InitCommon()
	pointscheme.InitPointScheme()
	memcache.InitMemCahce()
	his.InitHis()
	//history2.Init()
	httpr.InitHttpServer()
	//common.NotifyExit()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	cs := <-c
	//history2.Exit()
	logger.Infof("Got signal [%v], app exit.", cs)
}

package his

//处理历史db和存储的定时器
import (
	"github.com/robfig/cron"
	"github.com/zloathleo/industry_equipment_monitor/common/logger"
	"github.com/zloathleo/industry_equipment_monitor/memcache"
	"github.com/zloathleo/industry_equipment_monitor/db"
	"github.com/zloathleo/industry_equipment_monitor/utils"
	"fmt"
	"github.com/zloathleo/industry_equipment_monitor/systemalarm"
	"github.com/pkg/errors"
)

const (
	cornSpecCreateDBFile     = "0 59 23 * * ?"  //每天23:59整执行一次 -- 创建db文件
	cornSpecChangeCurrentDB  = "0 0 0 * * ?"    //每天0:0:0.500 执行一次 -- 切换db并插入新表缓存数据
	cornSpecInitNewTableInit = "0 0 1-23 * * ?" //每小时0:0:0.500 执行一次 0:0:0.500不执行 -- 插入新表缓存数据
	cornSpecUpdateHistory    = "* * * * * ?"    //每秒.950 执行一次 -- 存储历史数据
)

func initCorn() {
	c := cron.New()

	//每天23:59整执行一次 -- 创建db文件
	c.AddFunc(cornSpecCreateDBFile, func() {
		db.CreateTomorrowConnect()
	})

	//每天0:0:0.500 执行一次 -- 切换db并插入新表缓存数据
	c.AddFunc(cornSpecChangeCurrentDB, func() {
		exchangeTodayDB()
	})

	//每小时0:0:0.500 执行一次 0:0:0.500不执行 -- 插入新表缓存数据
	c.AddFunc(cornSpecInitNewTableInit, func() {
		initNewTable()

	})

	//每秒 整秒 执行一次 -- 存储历史数据
	c.AddFunc(cornSpecUpdateHistory, func() {
		saveRealtime(memcache.GlobalMemCache.CacheValueMap.SafeCopyAndClear(),false)
	})

	c.Start()
}

/**
切换DB文件
*/
func exchangeTodayDB() {
	connect := db.GetTodayConnect()

	if connect == nil {
		logger.Errorf("today db file %s is not exist", utils.GetTodayString())
		systemErr := errors.New(fmt.Sprintf("today db file %s is not exist", utils.GetTodayString()))
		systemalarm.AddSystemAlarm(systemErr)
	} else {
		todayDbLock.Lock()
		todayDb = connect.DB
		todayDbLock.Unlock()
		logger.Infof("exchange today db %s ok", utils.GetTodayString())
		initNewTable()
	}
}


/**
初始化当前表格,缓存中数据写入当前新表
启动
每小时0:0:0.500 执行
每天凌晨0:0:0.500执行
*/
func initNewTable() {
	logger.Info("begin insert memcache into newtable")
	list := memcache.GlobalMemCache.GetCurrentValueList()
	list["before-auto-copy"] = 99999999
	saveRealtime(list,true)
	logger.Info("finish insert memcache ")
}

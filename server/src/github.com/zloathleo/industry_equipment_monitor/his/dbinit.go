package his

import (
	"database/sql"
	"sync"
	"time"
	"github.com/zloathleo/industry_equipment_monitor/utils"
	"github.com/zloathleo/industry_equipment_monitor/common/logger"
	"os"
	"github.com/zloathleo/industry_equipment_monitor/db"
)

var (
	//当天的db 文件
	todayDb *sql.DB

	todayDbName string
	//当天db的读写锁
	todayDbLock = new(sync.RWMutex)
)


func initDB() {
	dir := "data"
	begin := time.Now().UnixNano()
	exist := utils.IsFileExist(dir)
	if exist {
		logger.Warnln("data file directory is exist.")
	} else {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			logger.Fatalf("db file directory can't create. %v", err)
		} else {
			logger.Warnln("db file directory created.")
		}
	}
	connect := db.CreateAnyConnect(time.Now())
	todayDb = connect.DB
	logger.Infof("init db const time %d ms.", (time.Now().UnixNano()-begin)/1000000)
}


package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"sync"
	"time"
	"github.com/zloathleo/industry_equipment_monitor/utils"
	"github.com/zloathleo/industry_equipment_monitor/common/logger"
	"github.com/zloathleo/industry_equipment_monitor/systemalarm"
	"github.com/zloathleo/industry_equipment_monitor/appconst"
)

const (
	his24table = `CREATE TABLE his24_%d (
    pn    VARCHAR (32) NOT NULL,
    t     INT (11)     NOT NULL,
    value DOUBLE       NOT NULL,
    PRIMARY KEY (
        pn,
        t
    )
	);`
)

var (
	//fileName -- > DB
	allDbMap = make(map[string]*SQLite3FileConnect, 32)
)

//根据时间获得当天的DB
func CreateAnyConnect(t time.Time) *SQLite3FileConnect {
	fileName := utils.GetDayString(t)
	dbConnect := allDbMap[fileName]
	if dbConnect != nil {
		return dbConnect
	} else {
		connect := NewSQLite3FileConnect(fileName, true)
		allDbMap[fileName] = connect
		return connect
	}
}

func CreateTomorrowConnect() *SQLite3FileConnect {
	tomorrow := utils.GetTimeNextDay0Hour(time.Now())
	return CreateAnyConnect(tomorrow)
}

//获得文件已经存在的连接,不创建新文件
func GetAnyConnect(t time.Time) *SQLite3FileConnect {
	fileName := utils.GetDayString(t)
	dbConnect := allDbMap[fileName]
	if dbConnect != nil {
		return dbConnect
	} else {
		connect := NewSQLite3FileConnect(fileName, false)
		allDbMap[fileName] = connect
		return connect
	}
}

func GetTodayConnect() *SQLite3FileConnect {
	return GetAnyConnect(time.Now())
}

func GetDBPath(fileName string) string {
	return "data/data-" + fileName + ".DB"
}

type SQLite3FileConnect struct {
	Name           string //fileName not path
	DBFilePathName string
	DB             *sql.DB
	DBLock         *sync.RWMutex
}

func NewSQLite3FileConnect(fileName string, create bool) *SQLite3FileConnect {
	dbFilePathName := GetDBPath(fileName)

	logger.Debugf("newSQLite3FileConnect %s ", dbFilePathName)
	connect := &SQLite3FileConnect{Name: fileName, DBFilePathName: dbFilePathName}
	if utils.IsFileExist(dbFilePathName) {
		logger.Infof("newSQLite3FileConnect %s is exist.", dbFilePathName)
		var err error
		connect.DB, err = connect.ConnectGivenFileDb()
		if err != nil || connect.DB == nil {
			logger.Errorf("connect given file DB %v err.", dbFilePathName)
		} else {
			connect.DBLock = new(sync.RWMutex)
			return connect
		}
	} else if create {
		_, err := os.Create(dbFilePathName)
		if err == nil {
			logger.Infof("newSQLite3FileConnect create %s file success.", dbFilePathName)
			connect.DB, err = connect.InitGivenFileDb()
			if err != nil || connect.DB == nil {
				logger.Errorf("init given file DB %v err.", dbFilePathName)
			} else {
				connect.DBLock = new(sync.RWMutex)
				return connect
			}
		} else {
			systemalarm.AddSystemAlarm(err)
			logger.Errorf("newSQLite3FileConnect create %s file is err %v.", dbFilePathName, err)
		}
	}
	return nil
}

/**
初始化指定db文件
*/
func (conn *SQLite3FileConnect) InitGivenFileDb() (*sql.DB, error) {
	cdb, err := conn.ConnectGivenFileDb()

	if cdb != nil && err == nil {
		tx, _ := cdb.Begin()
		for i := 0; i < 24; i++ {
			createTableSql := fmt.Sprintf(his24table, i)
			_, createTableErr := cdb.Exec(createTableSql)
			if createTableErr != nil {
				systemalarm.AddSystemAlarm(createTableErr)
			}
		}
		tx.Commit()
	}
	return cdb, err
}

/**
连接指定db文件
*/
func (conn *SQLite3FileConnect) ConnectGivenFileDb() (*sql.DB, error) {
	cdb, err := sql.Open(appconst.DriverName, conn.DBFilePathName)
	if err != nil {
		return nil, err
	}
	if err = cdb.Ping(); err != nil {
		return nil, err
	}
	return cdb, nil
}

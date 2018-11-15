package pointscheme

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zloathleo/industry_equipment_monitor/appconst"
	"github.com/zloathleo/industry_equipment_monitor/common/logger"
	"github.com/zloathleo/industry_equipment_monitor/db"
	"gopkg.in/Knetic/govaluate.v2"
)

//测点scheme 信息字典
var GlobalPointSchemeMap = NewConcurrentMap()

type Point struct {
	Name     string
	Type     int
	NickName *db.NullString
	Unit     *db.NullString
	Desc     *db.NullString

	Upper *db.NullFloat64
	Lower *db.NullFloat64
	H     *db.NullFloat64
	L     *db.NullFloat64
	HH    *db.NullFloat64
	LL    *db.NullFloat64

	ValueFormula           *db.NullString
	ValueFormulaExpression *govaluate.EvaluableExpression `json:"-"`
	AlarmFormula           *db.NullString
	AlarmFormulaExpression *govaluate.EvaluableExpression `json:"-"`
	Device                 *db.NullString
}

func InitPointScheme() {
	db, err := sql.Open(appconst.DriverName, appconst.SchemeDbFile)
	if err != nil {
		logger.Fatalf("scheme db file %s load err %v", appconst.SchemeDbFile, err)
	}
	if err = db.Ping(); err != nil {
		logger.Fatalf("scheme db file %s ping err %v", appconst.SchemeDbFile, err)
	}

	selectSql := "select t.name,t.typ,t.nickname,t.unit,t.des,t.upper,t.lower,t.h,t.l,t.hh,t.ll,t.value_formula,t.alarm_formula,t.device from " + appconst.SchemeTableName + " t"
	logger.Info(selectSql)
	rows, err := db.Query(selectSql)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if err != nil {
		logger.Fatalf("scheme db file '%s' query err %v", appconst.SchemeTableName, err)
	}

	for rows.Next() {
		//var col1, col2, col3, col4,col5, col6 interface{}
		var t Point
		err := rows.Scan(&t.Name, &t.Type, &t.NickName, &t.Unit, &t.Desc, &t.Upper, &t.Lower, &t.H, &t.L, &t.HH, &t.LL, &t.ValueFormula, &t.AlarmFormula, &t.Device)
		if err != nil {
			logger.Fatalf("scheme db file %s foreach err %v", appconst.SchemeTableName, err)
		} else {
			//初始化值公式
			if t.ValueFormula.Valid {
				valueFormulaStr := t.ValueFormula.String
				if len(valueFormulaStr) > 0 {
					expression, err := govaluate.NewEvaluableExpression(valueFormulaStr)
					if err != nil {
						logger.Errorf("point '%s' scheme value_formula '%s' is error.", t.Name, valueFormulaStr)
					} else {
						t.ValueFormulaExpression = expression
					}
				}
			}

			//初始化报警公式
			if t.AlarmFormula.Valid {
				alarmFormulaStr := t.AlarmFormula.String
				if len(alarmFormulaStr) > 0 {
					expression, err := govaluate.NewEvaluableExpression(alarmFormulaStr)
					if err != nil {
						logger.Errorf("point '%s' scheme alarm_formula '%s' is error.", t.Name, alarmFormulaStr)
					} else {
						t.AlarmFormulaExpression = expression
					}
				}
			}
			GlobalPointSchemeMap.Set(t.Name,&t)
		}
	}

	logger.Warnln("pointscheme init ok.")
}

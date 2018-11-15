package realtime

import (
	"github.com/zloathleo/industry_equipment_monitor/pointscheme"
	"math"
	"github.com/zloathleo/industry_equipment_monitor/common/logger"
)

//值转换
func translatePointValue(key string, value float64) float64 {
	pointScheme,exist := pointscheme.GlobalPointSchemeMap.Get(key)
	if pointScheme != nil && exist{
		conversion := pointScheme.ValueFormulaExpression
		if conversion != nil {
			parameters := make(map[string]interface{}, 8)
			parameters["x"] = value
			result, err := conversion.Evaluate(parameters)
			if err == nil {
				value = result.(float64)
			} else {
				logger.Infof("scheme conversion %v for %v is error", conversion.String(), value)
			}
		}
	}else{
		//未配置的点
		logger.Debugf("point '%s' scheme is not found.",key)
	}
	if math.Abs(value-0) < 0.00001 {
		value = 0
	}
	return value
}

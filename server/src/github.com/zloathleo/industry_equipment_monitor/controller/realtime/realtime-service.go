package realtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zloathleo/industry_equipment_monitor/dstruct"
	"github.com/zloathleo/industry_equipment_monitor/memcache"
	"strconv"
	"strings"
)

func generatePointRealtimeValue(points string) ([]byte, error) {
	if points == "" {
		return nil, errors.New(fmt.Sprintf("param 'names' is '%s'.", points))
	}
	pointsArray := strings.Split(points, ",")
	if pointsArray == nil {
		return nil, errors.New(fmt.Sprintf("param 'names' [%s] is err.", points))
	} else {

		var builder bytes.Buffer

		//root
		builder.WriteString("{ ")
		for _, pn := range pointsArray {
			builder.WriteString("\"" + pn + "\":")
			exits, value := memcache.GlobalMemCache.GetCurrentValue(pn)
			if exits {
				builder.WriteString(strconv.FormatFloat(value, 'f', 2, 64) + ",")
			} else {
				builder.WriteString("null,")
			}
		}
		if builder.Len() > 4 {
			builder.Truncate(builder.Len() - 1)
		}
		//root
		builder.WriteString(" }")

		return builder.Bytes(), nil

	}
}

func pushRealtimeValue(dataJson string) ([]byte, error) {
	var pushDas dstruct.PushDas
	err := json.Unmarshal([]byte(dataJson), &pushDas)
	if err == nil {
		rows := pushDas.Rows
		for _, item := range rows {
			tv := translatePointValue(item.PointName, item.Value)
			memcache.GlobalMemCache.SaveCurrentValue(item.PointName, tv)

		}
		return nil, nil
	} else {
		return nil, err
	}
}

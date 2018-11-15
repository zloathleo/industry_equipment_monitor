package scheme

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin/json"
	"github.com/zloathleo/industry_equipment_monitor/pointscheme"
)

//生成点信息
func generatePointSchemeMap() ([]byte, error) {
	return json.Marshal(pointscheme.GlobalPointSchemeMap)
}

func generatePointScheme(name string) ([]byte, error) {
	point, exist := pointscheme.GlobalPointSchemeMap.Get(name)
	if point == nil || !exist {
		return nil, errors.New(fmt.Sprintf("point '%s' scheme is not found.", name))
	} else {
		return json.Marshal(point)
	}
}

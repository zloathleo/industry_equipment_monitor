package history

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

func generatePointHistoryChart(chartType string, points string, to string, dur string, interval string) ([]byte, error) {
	if points == "" {
		return nil, errors.New(fmt.Sprintf("param 'points' is '%s'.", points))
	}
	pointsArray := strings.Split(points, ",")
	if pointsArray == nil || len(pointsArray) == 0{
		return nil, errors.New(fmt.Sprintf("param 'points' [%s] is err.", points))
	}

	to64, err := strconv.ParseInt(to, 10, 64)
	if err != nil {
		return nil, err
	}
	durInt, err := strconv.Atoi(dur)
	if err != nil {
		return nil, err
	}
	intervalInt := 0
	if interval != "" {
		intervalInt, err = strconv.Atoi(interval)
		if err != nil {
			return nil, err
		}
	}

	historyMap := fetchHistoryChartData(pointsArray, to64, durInt, intervalInt)
	if historyMap != nil {

		jsonBuffer := renderChartHistoryJson(historyMap)
		return jsonBuffer.Bytes(), nil

		//if chartType == "radar" {
		//	//jsonBuffer := renderRadarChartHistoryJson(pointsArray, hisMap, xAxis)
		//	//return jsonBuffer.Bytes(), nil
		//	//c.Status(http.StatusOK)
		//	//c.Writer.Write(jsonBuffer.Bytes())
		//} else {
		//	jsonBuffer := renderChartHistoryJson(hisMap, xAxis)
		//	return jsonBuffer.Bytes(), nil
		//	//c.Status(http.StatusOK)
		//	//c.Writer.Write(jsonBuffer.Bytes())
		//}
	} else {
		return []byte("{}"), nil
	}

}

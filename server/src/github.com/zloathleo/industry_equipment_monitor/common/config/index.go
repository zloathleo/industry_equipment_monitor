package config

import (
	"github.com/zloathleo/industry_equipment_monitor/common/logger"
)

const LogoAsic = ` 
***********************************************************************************
****                        	   SAFEFIRE INC.                               ****
*********************************************************************************** 
**** load config file ok
**** App.Name %s
**** App.Port %d
**** Control.AutoLogic %v
***********************************************************************************
`

var AppConfig = struct {
	App struct {
		Name     string `default:"center-server"`
		Port     uint   `required:"true" default:"8088"`
		LogLevel string `default:"W"`
	}
	Control struct {
		AutoLogic bool   `default:false`
		DasHost      string `yaml:"das-host" default:"localhost"`
		DasPort      uint   `yaml:"das-port" default:"8080"`
	}
}{}

func Init() {
	err := Load(&AppConfig, "config.yml")
	if err != nil {
		logger.Fatalf("load config file error [%v]", err)
	} else {
		logger.Infof(LogoAsic, AppConfig.App.Name, AppConfig.App.Port, AppConfig.Control.AutoLogic)
	}
}

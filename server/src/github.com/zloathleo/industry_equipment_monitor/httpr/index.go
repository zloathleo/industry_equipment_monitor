package httpr

import (
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/zloathleo/industry_equipment_monitor/common/config"
	"github.com/zloathleo/industry_equipment_monitor/common/logger"
	"github.com/zloathleo/industry_equipment_monitor/controller"
	"net/http"
	"time"
)

func InitHttpServer() {

	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = logger.GetOutput()

	router := gin.New()
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Static("/public", "./public")

	router.Use(newCrosMiddleware())
	dasGroup := router.Group("/")
	controller.InitController(dasGroup)

	addr := fmt.Sprintf(":%d", config.AppConfig.App.Port)
	s := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go s.ListenAndServe()
	logger.Warnln("http init ok.")
	logger.Warnf("http server is running and listen at %d.", config.AppConfig.App.Port)
}

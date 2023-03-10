package api

import (
	"context"
	"fmt"
	"github.com/auribuo/novasearch/api/features"
	"github.com/auribuo/novasearch/api/middleware"
	"github.com/auribuo/novasearch/log"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Start(host string, port int, debug bool) {
	gin.SetMode(gin.ReleaseMode)
	if debug {
		gin.SetMode(gin.DebugMode)
	}
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(log.GinLogger())
	router.Use(middleware.CORS())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	var api = router.Group("/api")
	api.GET("/health", func(context *gin.Context) {
		context.Status(200)
	})

	api.GET("/galaxies", features.AllGalaxies)
	api.GET("/galaxies/:id", features.Galaxy)
	api.POST("/galaxies", features.FilterGalaxies)
	api.GET("/calculate/galaxies", features.CalculateGalaxyPath)
	api.GET("/calculate/viewports", features.CalculateViewportPath)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: router,
	}

	processSignalChan := make(chan os.Signal, 1)
	signal.Notify(processSignalChan, os.Interrupt, os.Kill)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Logger.Error(err.Error())
		}
	}()

	<-processSignalChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Logger.Error(err.Error())
	}
}

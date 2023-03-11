package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/auribuo/novasearch/api/middleware"
	"github.com/auribuo/novasearch/api/routes"
	"github.com/auribuo/novasearch/log"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Start(host string, port int, devMode bool) {
	gin.SetMode(gin.ReleaseMode)
	if devMode {
		gin.SetMode(gin.DebugMode)
	}
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(log.GinMiddleware())
	router.Use(middleware.CORS())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	var api = router.Group("/api")
	api.GET("/health", func(context *gin.Context) {
		context.Status(200)
	})

	api.GET("/galaxies", routes.AllGalaxies)
	api.GET("/galaxies/:id", routes.Galaxy)
	api.POST("/galaxies", routes.FilterGalaxies)
	api.GET("/calculate/galaxies", routes.CalculateGalaxyPath)
	api.GET("/calculate/viewports", routes.CalculateViewportPath)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: router,
	}

	processSignalChan := make(chan os.Signal, 1)
	signal.Notify(processSignalChan, os.Interrupt)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Default.Error(err.Error())
			} else {
				log.Default.Info("server closed")
			}
		}
	}()

	log.Default.Info(fmt.Sprintf("server listening on http://%s:%d", host, port))

	<-processSignalChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Default.Error(err.Error())
	}
}

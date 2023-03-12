package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/auribuo/novasearch/api/middleware"
	"github.com/auribuo/novasearch/docs"
	"github.com/auribuo/novasearch/log"
	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Start(host string, port int, devMode bool) {
	docs.SwaggerInfo.BasePath = "/api"
	gin.SetMode(gin.ReleaseMode)
	if devMode {
		gin.SetMode(gin.DebugMode)
	}

	memoryStore := persist.NewMemoryStore(5 * time.Minute)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(log.GinMiddleware())
	router.Use(middleware.CORS())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/", func(context *gin.Context) {
		context.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	var api = router.Group("/api")
	api.GET("/health", func(context *gin.Context) {
		context.Status(200)
	})

	api.GET("/galaxies", cache.CacheByRequestURI(memoryStore, 1*time.Minute), AllGalaxies)
	api.GET("/galaxies/:id", Galaxy)
	api.POST("/galaxies", FilterGalaxies)
	api.POST("/calculate/galaxies/:strategy", CalculateGalaxyPath)
	api.POST("/calculate/viewports/:strategy", CalculateViewportPath)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: router,
	}

	processSignalChan := make(chan os.Signal, 1)
	signal.Notify(processSignalChan, os.Interrupt)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Default.Error(err)
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
		log.Default.Error(err)
	}
}

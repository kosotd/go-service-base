package server

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kosotd/go-service-base/cache"
	"github.com/kosotd/go-service-base/config"
	"github.com/kosotd/go-service-base/database"
	"github.com/kosotd/go-service-base/utils"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var router *gin.Engine

func init() {
	if utils.Equals(config.BuildMode(), gin.ReleaseMode) {
		gin.SetMode(gin.ReleaseMode)
	} else if utils.Equals(config.BuildMode(), gin.DebugMode) {
		gin.SetMode(gin.DebugMode)
	}

	router = gin.New()

	router.Use(cors.New(cors.Config{
		AllowHeaders: allowHeaders,
		AllowMethods: allowMethods,
		AllowOrigins: config.AllowedOrigins(),
	}))
}

func AddHandler(method string, path string, handler http.HandlerFunc) {
	router.Handle(method, path, gin.WrapF(handler))
}

func AddGetHandler(path string, handler http.HandlerFunc) {
	router.GET(path, gin.WrapF(handler))
}

func AddPostHandler(path string, handler http.HandlerFunc) {
	router.POST(path, gin.WrapF(handler))
}

func AddGinHandler(method string, path string, handler gin.HandlerFunc) {
	router.Handle(method, path, handler)
}

func AddGinGetHandler(path string, handler gin.HandlerFunc) {
	router.GET(path, handler)
}

func AddGinPostHandler(path string, handler gin.HandlerFunc) {
	router.POST(path, handler)
}

func GetRouter() *gin.Engine {
	return router
}

func RunServer() {
	defer cache.Close()
	defer database.Close()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.ServerPort()),
		Handler: router,
	}

	if readTimeout, err := config.Duration("ReadTimeout"); err == nil {
		srv.ReadTimeout = readTimeout
	}

	if writeTimeout, err := config.Duration("WriteTimeout"); err == nil {
		srv.WriteTimeout = writeTimeout
	}

	crt, crtErr := config.String("CrtFile")
	key, keyErr := config.String("KeyFile")

	go func() {
		utils.LogInfo(fmt.Sprintf("server started on port: %s", config.ServerPort()))

		if crtErr == nil && keyErr == nil {
			if err := srv.ListenAndServeTLS(crt, key); err != nil && err != http.ErrServerClosed {
				utils.FailIfError(errors.Wrapf(err, "error start server"))
			}
		} else {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				utils.FailIfError(errors.Wrapf(err, "error start server"))
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	<-quit
	utils.LogInfo("shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		utils.FailIfError(errors.Wrap(err, "error shutdown server"))
	}

	utils.LogInfo("server exiting")
}

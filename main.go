package main

import (
	"example.com/vesssels-api/internal/vesselsapi"
	"example.com/vesssels-api/pkg/base"
	"github.com/gin-gonic/gin"
)

type Config struct {
	API base.Heartbeat
	Log base.Log
}

func main() {
	log := base.GetLogger()
	log.Info("Hello, world!")
	defer base.PanicHandler()
	c := base.GetConfig(&Config{})
	base.SetupLog(c.Log)

	store := vesselsapi.NewInMemoryStore()

	handler := vesselsapi.Handler{
		Log:   log,
		Store: store,
	}

	// register all routes with handlers
	router := gin.New()

	router.GET("/heartbeat", gin.WrapF(base.HeartbeatHandler))
	router.GET("/v1/vessels/:imo", handler.GetVesselByIMO)

	// start serving requests
	log.Fatal(router.Run(c.API.Addr))
}

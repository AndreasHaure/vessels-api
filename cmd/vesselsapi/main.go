package main

import (
	"example.com/vesssels-api/internal/vesselsapi"
	"example.com/vesssels-api/pkg/base"
	"github.com/gin-gonic/gin"
)

type Config struct {
	API      base.Heartbeat
	Log      base.Log
	Postgres base.Postgres
}

func main() {
	log := base.GetLogger()
	defer base.PanicHandler()
	c := base.GetConfig(&Config{})
	base.SetupLog(c.Log)
	db := base.SetupPostgres(c.Postgres)
	defer db.Close()

	store := vesselsapi.NewPGStore(db, c.Postgres.SchemaName)

	handler := vesselsapi.Handler{
		Log:   log,
		Store: store,
	}

	// register all routes with handlers
	router := gin.New()

	router.GET("/heartbeat", gin.WrapF(base.HeartbeatHandler))
	router.GET("/v1/vessels", handler.GetVessels)
	router.PUT("/v1/vessels/:imo", handler.UpdateVessel)
	router.GET("/v1/vessels/:imo", handler.GetVesselByIMO)
	router.DELETE("/v1/vessels/:imo", handler.DeleteVessel)

	// start serving requests
	log.Fatal(router.Run(c.API.Addr))
}

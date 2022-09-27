// Package classification keyplay-metadata-api API
//
// keyplay-metadata-api
//
//     Schemes: http
//     Host: services-dev.fubo.tv
//     BasePath: /keyplay-metadata-api
//     Version: 1.0.0
//
//     Produces:
//     - application/json
//
//
// swagger:meta
package main

import (
	"fmt"
	"log"

	"github.com/fubotv/fubotv-logging/v3/logging"
	"github.com/fubotv/fubotv-metrics/v12/datadog/spans"
	"github.com/fubotv/keyplay-metadata-api/app/config"
	"github.com/fubotv/keyplay-metadata-api/app/db"
	"github.com/fubotv/keyplay-metadata-api/app/server"
	"github.com/fubotv/keyplay-metadata-api/app/server/handler"
)

func main() {
	logging.InfoWithoutContext("keyplay-metadata-api launches")

	cfg := config.New()

	// create database client
	dbclient, err := db.CreateDBHandler()
	if err != nil {
		log.Fatal(fmt.Errorf("err creating bigtable client : %v", err))
	}

	keyplayHandler := handler.ServiceHandler{DatabaseHandler: dbclient}

	defer dbclient.Client.Close()
	// init tracer
	err = spans.InitTracer(cfg.DataDogCfg, cfg.AppCfg.ServiceName, cfg.AppCfg.Env, cfg.AppCfg.Version.BuildSource)
	if err != nil {
		log.Fatal(err)
	}

	// start server
	srv := server.New(cfg, keyplayHandler)
	srv.Run()
}

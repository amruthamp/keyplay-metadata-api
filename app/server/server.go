package server

import (
	"compress/gzip"
	"context"
	"fmt"

	"github.com/NYTimes/gziphandler"
	"github.com/fubotv/fubotv-http-utils/v21/middlewr"
	"github.com/fubotv/fubotv-logging/v3/logging"
	"github.com/fubotv/keyplay-metadata-api/app/config"
	"github.com/fubotv/keyplay-metadata-api/app/server/handler"
	"github.com/fubotv/keyplay-metadata-api/app/server/routes"

	"net/http"

	"goji.io"
	"goji.io/pat"
)

type HttpServer struct {
	cfg *config.Cfg
	mux *goji.Mux
}

func New(cfg *config.Cfg, serviceHandler handler.ServiceHandler) *HttpServer {
	s := &HttpServer{
		cfg: cfg,
		mux: goji.NewMux(),
	}

	//configure handlers
	s.mux.HandleFunc(pat.Get(routes.Healthcheck), handler.HealthcheckStatus)
	s.mux.HandleFunc(pat.Get(routes.SampleRoute), handler.SampleHandler)
	s.mux.HandleFunc(pat.Get(routes.GetKeyplayMetadata), serviceHandler.GetKeyplayMetadata)
	s.mux.HandleFunc(pat.Post(routes.CreateKeyplayMetadata), serviceHandler.CreateKeyPlayMetadata)
	s.mux.HandleFunc(pat.Delete(routes.DeleteKeyplayMetadata), serviceHandler.DeleteKeyplayMetadata)

	if s.cfg.AppCfg.Env == "local" || s.cfg.AppCfg.Env == "development" {
		s.mux.Handle(pat.Get(routes.Swagger+"/*"), http.StripPrefix(routes.Swagger+"/", http.FileServer(http.Dir(s.cfg.AppCfg.SwaggerDir))))
	}

	allRoutes := goji.SubMux()
	allRoutes.Use(middlewr.SendMetrics(cfg.AppCfg.ServiceName, []string{}))

	gzipWrapper, _ := gziphandler.NewGzipLevelHandler(gzip.BestSpeed)
	s.mux.Use(gzipWrapper)

	return s
}

func (s *HttpServer) Run() {
	address := fmt.Sprintf(":%d", s.cfg.AppCfg.ServerPort)

	logging.Info(context.Background(), "Starting Http Server on port ", s.cfg.AppCfg.ServerPort)
	if err := http.ListenAndServe(address, s.mux); err != nil {
		panic(err)
	}
}

func (s *HttpServer) GetMux() *goji.Mux {
	return s.mux
}

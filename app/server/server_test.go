package server_test

// import (
// 	"github.com/fubotv/keyplay-metadata-api/app/config"
// 	. "github.com/fubotv/keyplay-metadata-api/app/server"
// 	"github.com/fubotv/keyplay-metadata-api/app/server/routes"
// 	. "github.com/smartystreets/goconvey/convey"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"testing"
// )

// func TestMain(m *testing.M) {
// 	os.Exit(m.Run())
// }

// func TestHttpServer_Runs(t *testing.T) {
// 	srv := New(config.GetConfig())
// 	go srv.Run()

// 	Convey("http server test", t, func() {
// 		Convey("provision HealthcheckStatus endpoint", func() {
// 			req, err := http.NewRequest(http.MethodGet, routes.Healthcheck, nil)
// 			So(err, ShouldBeNil)
// 			rr := httptest.NewRecorder()
// 			srv.GetMux().ServeHTTP(rr, req)
// 			So(rr.Code, ShouldEqual, http.StatusOK)
// 		})
// 	})
// }

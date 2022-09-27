package handler_test

import (
	"encoding/json"
	"github.com/fubotv/keyplay-metadata-api/app/server/handler"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthcheckStatus(t *testing.T) {
	Convey("Test health check handler", t, func() {
		req, err := http.NewRequest("GET", "/healthcheck", nil)
		So(err, ShouldBeNil)
		rr := httptest.NewRecorder()
		h := http.HandlerFunc(handler.HealthcheckStatus)
		h.ServeHTTP(rr, req)
		So(rr.Code, ShouldEqual, http.StatusOK)
		var statuses []handler.Status
		decoder := json.NewDecoder(rr.Body)
		err = decoder.Decode(&statuses)
		So(err, ShouldBeNil)
	})
}

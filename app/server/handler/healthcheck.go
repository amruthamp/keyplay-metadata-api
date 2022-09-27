package handler

import (
	"net/http"

	"github.com/fubotv/fubotv-http-utils/v21/respx"
)

// var empty interface{}

func HealthcheckStatus(writer http.ResponseWriter, request *http.Request) {
	var statusList []Status
	currentStatus := http.StatusOK
	statusList = append(statusList, Status{
		Status:  200,
		Name:    "keyplay-metadata-api",
		Message: "OK",
	})

	respx.Write(request.Context(), writer, nil, currentStatus, &statusList)
}

type Status struct {
	Name    string `json:"name"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

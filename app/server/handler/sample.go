package handler

import (
	"github.com/fubotv/fubotv-http-utils/v21/respx"
	"net/http"
)

// swagger:model response
type Response struct {
	Data interface{} `json:"data"`
}

// swagger:model error
type ErrorDTO struct {
	Error ErrorDetails `json:"error"`
}

// swagger:model errorMessage
type ErrorDetails struct {
	Message string `json:"message"`
}

// swagger:route GET /sample GetSampleV1
//
// Responses:
// 		    200: response
//			400: error
//			500: error
func SampleHandler(writer http.ResponseWriter, request *http.Request) {
	respx.Write(request.Context(), writer, nil, http.StatusOK,  &Response{Data: "sample response"})
}
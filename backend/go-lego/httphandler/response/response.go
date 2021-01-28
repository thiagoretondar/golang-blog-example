// Package response contains some method for handling http.ResponseWriter
package response

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/middleware"
)

type HTTPResponse struct {
	Data  interface{}  `json:"data,omitempty"`
	Error *ClientError `json:"error,omitempty"`
}

type ClientError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// WithJSON writes a http response with specific code status and JSON
func WithJSON(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	// Set content type header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Set response header with X-request-ID
	w.Header().Set("X-request-ID", middleware.GetReqID(r.Context()))

	// Write status
	w.WriteHeader(code)

	// Write response
	if code != http.StatusNoContent && data != nil {
		jsonResponse, err := json.Marshal(data)
		// Marshal response as json
		if err != nil {
			WithJSONError(w, r, http.StatusInternalServerError, err)
			return
		}
		_, err = w.Write(jsonResponse)
		if err != nil {
			WithJSONError(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

// WithJSONError writes a http response with specific code status and JSON with "error" field
func WithJSONError(w http.ResponseWriter, r *http.Request, code int, err error) {
	if err != nil {
		WithJSON(w, r, code, &ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	} else {
		WithJSON(w, r, code, &ErrorResponse{
			Error:   true,
			Message: http.StatusText(code),
		})
	}
}

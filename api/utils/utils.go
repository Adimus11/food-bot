package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ApiResponse struct {
	StatusCode int
	Response   json.RawMessage
}

type Request interface {
	Payload() interface{}
	Do(ctx context.Context, payload interface{}) (*ApiResponse, error)
}

func SendJSON(resp *ApiResponse, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(resp.Response)
}

func SendError(err error, w http.ResponseWriter) {
	fmt.Println(err.Error())
	w.Header().Set("Content-Type", "application/json")
	if apiError, ok := err.(*ApiError); ok {
		if marshaledError, err := json.Marshal(apiError); err == nil {
			w.WriteHeader(apiError.StatusCode)
			w.Write(marshaledError)
			return
		}
	}
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal server error"))
}

func HandlerWrapper(req Request) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := req.Payload()

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method == "POST" || r.Method == "PUT" {
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&payload)
			if err != nil {
				SendError(err, w)
				return
			}
		}

		resp, err := req.Do(r.Context(), payload)
		if err != nil {
			SendError(err, w)
			return
		}
		SendJSON(resp, w)
	}
}

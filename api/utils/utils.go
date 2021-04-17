package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiResponse struct {
	StatusCode int
	Response   interface{}
	Cookies    []*http.Cookie
}

type Request interface {
	Payload() interface{}
	Do(ctx context.Context, vars map[string]string, payload interface{}) (*ApiResponse, error)
}

func SendJSON(resp *ApiResponse, w http.ResponseWriter) {
	bytesResponse, err := json.Marshal(resp.Response)
	if err != nil {
		SendError(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	for _, cookie := range resp.Cookies {
		http.SetCookie(w, cookie)
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(bytesResponse)
}

func SendError(err error, w http.ResponseWriter) {
	fmt.Printf("Error during processing: %s\n", err)
	w.Header().Set("Content-Type", "application/json")
	if apiError, ok := err.(*ApiError); ok {
		if marshaledError, err := json.Marshal(apiError); err == nil {
			w.WriteHeader(apiError.StatusCode)
			w.Write(marshaledError)
			return
		}
	}
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("{ \"error\": \"internal server error\" }"))
}

func HandlerWrapper(req Request) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := req.Payload()

		if r.Method == "OPTIONS" {
			w.Write([]byte(""))
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodDelete && r.Method != http.MethodGet && r.Method != http.MethodOptions {
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&payload)
			if err != nil {
				SendError(err, w)
				return
			}
		}

		resp, err := req.Do(r.Context(), mux.Vars(r), payload)
		if err != nil {
			SendError(err, w)
			return
		}
		SendJSON(resp, w)
	}
}

package main

import (
    "encoding/json"
    "net/http"
)

type helloResp struct {
    Message string `json:"message"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    resp := helloResp{Message: "Hello from Cloud Run-ready Go service!"}
    _ = json.NewEncoder(w).Encode(resp)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte("ok"))
}

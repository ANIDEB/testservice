package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHelloHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/", nil)
    rr := httptest.NewRecorder()

    helloHandler(rr, req)

    if rr.Code != http.StatusOK {
        t.Fatalf("expected status 200; got %d", rr.Code)
    }

    expected := `{"message":"Hello from Cloud Run-ready Go service!"}`
    got := rr.Body.String()
    // JSON encoder appends a newline; trim that by simple comparison allowance
    if got != expected+"\n" && got != expected {
        t.Fatalf("unexpected body: %q", got)
    }
}

func TestHealthHandler(t *testing.T) {
    paths := []string{"/healthz", "/health", "/readyz", "/livez", "/_ah/health"}
    for _, p := range paths {
        req := httptest.NewRequest("GET", p, nil)
        rr := httptest.NewRecorder()

        healthHandler(rr, req)

        if rr.Code != http.StatusOK {
            t.Fatalf("%s: expected status 200; got %d", p, rr.Code)
        }
        if rr.Body.String() != "ok" {
            t.Fatalf("%s: expected body 'ok', got %q", p, rr.Body.String())
        }
    }
}

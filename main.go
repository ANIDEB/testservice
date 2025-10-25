package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", helloHandler)
    // Health endpoints (Cloud Run may not forward /healthz externally)
    mux.HandleFunc("/healthz", healthHandler)
    mux.HandleFunc("/health", healthHandler)
    mux.HandleFunc("/readyz", healthHandler)
    mux.HandleFunc("/livez", healthHandler)
    mux.HandleFunc("/_ah/health", healthHandler)

    srv := &http.Server{
        Addr:    ":" + getPort(),
        Handler: mux,
    }

    // start server
    go func() {
        log.Printf("starting server on %s", srv.Addr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %s\n", err)
        }
    }()

    // wait for interrupt
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }
    log.Println("server stopped")
}

func getPort() string {
    p := os.Getenv("PORT")
    if p == "" {
        p = "8080"
    }
    return p
}

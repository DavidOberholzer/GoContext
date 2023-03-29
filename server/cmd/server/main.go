package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"server/internal/handler"
	"server/internal/storage"
)

func main() {
	// This context is not in this example, but imagine there were some other processes
	// that could handle cancellation for graceful shutdown.
	_, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	store := storage.NewStore()
	h := handler.New(store)

	server := http.Server{
		Addr:        ":8080",
		Handler:     h,
		ReadTimeout: time.Second * 10,
	}

	wg.Add(1)
	// Start server in go routine to avoid blocking here.
	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println("error with server")
		}
	}()

	shutdown := func() {
		defer wg.Done()
		// Timeout shutdown if taking too long.
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			fmt.Printf("http server shutdown error: %v", err)
		}
		fmt.Println("http server shutdown complete")
	}
	fmt.Println("TCP server listening on port :8080...")

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	<-sigs

	// Shutdown HTTP server and then cancel context for others using it.
	shutdown()
	cancel()
	// Wait for all routines to finish up before shutdown.
	wg.Wait()
	fmt.Println("Server application shutdown")
}

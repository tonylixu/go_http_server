// GO HTTP Server
package main

// Import necessary packages
import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"./handlers"

	"net/http/pprof"
	_ "net/http/pprof"

	"github.com/gorilla/mux" // need to use dep for package management
)

// HTTP server debug request handlers
func AttachProfiler(router *mux.Router) {
	router.HandleFunc("/debug/pprof", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
}

// main function
func main() {
	// Implements a request router and dispatcher for matching
	// incoming requests to their respective handler.
	router := mux.NewRouter()
	AttachProfiler(router)

	// Register handlers
	router.HandleFunc("/", handlers.RootHandler).Methods("GET")
	router.HandleFunc("/headers", handlers.HeaderHandler).Methods("GET")
	router.HandleFunc("/healthz", handlers.SuccessHandler).Methods("GET")

	// Create new http.Server object
	srv := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Create a new channel of capacity 1
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start a goroutine for http server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server started")

	<-done
	fmt.Println("Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}

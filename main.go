// GO HTTP Server
package main

// Import necessary packages
import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	handlers "github.com/tonylixu/go_http_server/handlers"
	probe "github.com/tonylixu/go_http_server/probe"

	"net/http/pprof"
	_ "net/http/pprof"

	"github.com/gorilla/mux" // need to use dep for package management
	"go.uber.org/zap"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tonylixu/go_http_server/metrics"
)

// HTTP server debug request handlers
func AttachProfiler(router *mux.Router) {
	router.HandleFunc("/debug/pprof", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
}

func ParseArguments() (int, string) {
	httpServerPort := flag.Int("port", 8080, "HTTP server port")
	logFile := flag.String("log", "/logs/http_server.log", "Log file location")

	flag.Parse()
	return *httpServerPort, *logFile
}

// main function
func main() {
	// Parse command line arguments
	httpServerPort, logFile := ParseArguments()
	fmt.Println("Port is ", httpServerPort, "log is ", logFile)
	metrics.Register()
	// Implements a request router and dispatcher for matching
	// incoming requests to their respective handler.
	router := mux.NewRouter()
	AttachProfiler(router)

	// Register handlers
	router.HandleFunc("/", handlers.RootHandler).Methods("GET")
	router.HandleFunc("/headers", handlers.HeaderHandler).Methods("GET")
	router.HandleFunc("/healthz", handlers.SuccessHandler).Methods("GET")
	router.Path("/metrics").Handler(promhttp.Handler())

	//Create a configuration for logger
	config := zap.NewProductionConfig()
	// config.OutputPaths = []string{"/logs/http-server.log"}
	config.OutputPaths = []string{logFile}
	zapLogger, err := config.Build()
	if err != nil {
		zapLogger.Error(fmt.Sprint("Can't initialize zap logger", err))
	}

	// Create new http.Server object
	httpServerPortString := fmt.Sprintf(":%d", httpServerPort)
	srv := http.Server{
		Addr:    httpServerPortString,
		Handler: router,
	}

	// Create a new channel of capacity 1
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start a goroutine for http server
	go func() {
		if err := probe.Create(); err != nil {
			panic(err)
		}

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Error(fmt.Sprint("Error listening: ", err))
		}

		if err := probe.Remove(); err != nil {
			panic(err)
		}
	}()
	zapLogger.Info("Server started")

	<-done
	zapLogger.Info("Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		zapLogger.Error(fmt.Sprint("Server Shutdown Failed: ", err))
	}

	// Flushes buffer
	defer zapLogger.Sync()
	zapLogger.Error("Server Exited Properly")
}

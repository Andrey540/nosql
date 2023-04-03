package server

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

// ServeFunc - runs server
type ServeFunc func() error

// StopFunc - stops server
type StopFunc func() error

const (
	serverIsCreated int32 = iota
	serverIsRunning
	serverIsStopped
)

type server struct {
	serveFunc ServeFunc
	stopFunc  StopFunc
	state     int32
}

func newServer(serve ServeFunc, stop StopFunc) *server {
	return &server{
		serveFunc: serve,
		stopFunc:  stop,
		state:     serverIsCreated,
	}
}

func (s *server) serve() error {
	if !atomic.CompareAndSwapInt32(&s.state, serverIsCreated, serverIsRunning) {
		if atomic.LoadInt32(&s.state) == serverIsRunning {
			return errAlreadyRun
		}
		return errTryRunStoppedServer
	}
	return s.serveFunc()
}

func (s *server) stop() error {
	stopped := atomic.CompareAndSwapInt32(&s.state, serverIsCreated, hubIsStopped) ||
		atomic.CompareAndSwapInt32(&s.state, serverIsRunning, serverIsStopped)

	if !stopped {
		return errAlreadyStopped
	}
	return s.stopFunc()
}

func ListenOSKillSignals(stopChan chan<- struct{}) {
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
		<-ch
		stopChan <- struct{}{}
	}()
}

func ServeHTTP(
	serveRESTAddress string,
	serverHub *Hub,
	handlerEntity http.HandlerFunc,
	handlerName http.HandlerFunc,
	logger, errorLogger *stdlog.Logger,
) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	var httpServer *http.Server
	serverHub.Serve(func() error {
		router := mux.NewRouter()
		router.HandleFunc("/entity", handlerEntity).Methods(http.MethodGet)
		router.HandleFunc("/name", handlerName).Methods(http.MethodGet)

		router.Use(RecoverMiddleware(errorLogger))

		httpServer = &http.Server{
			Handler:      router,
			Addr:         serveRESTAddress,
			ReadTimeout:  time.Hour,
			WriteTimeout: time.Hour,
		}

		return httpServer.ListenAndServe()
	}, func() error {
		cancel()
		return httpServer.Shutdown(context.Background())
	})
}

func InitLogger() *stdlog.Logger {
	return stdlog.New(os.Stdout, "http: ", stdlog.LstdFlags)
}

func InitErrorLogger() *stdlog.Logger {
	return stdlog.New(os.Stderr, "http: ", stdlog.LstdFlags)
}

var errAlreadyStopped = errors.New("server is not running, can't change server state")
var errAlreadyRun = errors.New("server is running, can't change server state to running")
var errTryRunStoppedServer = errors.New("server is stopped, can't change server state to running")

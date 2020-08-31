package server

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type Endpoint struct {
	Path        string
	Verb        string
	HandlerFunc http.HandlerFunc
}

func Start(port int, endpoints ...Endpoint) error {
	r := mux.NewRouter()
	for _, v := range endpoints {
		r.HandleFunc(v.Path, v.HandlerFunc).Methods(v.Verb)
	}
	http.Handle("/", r)

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%d", port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	shutDownC := make(chan error)
	listenC := make(chan error)

	go func() {
		log.Info("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			listenC <- err
		}
	}()

	go waitForShutdown(srv, shutDownC)

	select {
	case err := <-shutDownC:
		return err
	case err := <-listenC:
		return err
	}
}

func waitForShutdown(srv *http.Server, errChan chan<- error) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	sign := <-interruptChan

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		errChan <- fmt.Errorf("os signal %v received, err when shutdown: %v\n", sign, err)
		return
	}
	errChan <- fmt.Errorf("os signal %v received\n", sign)
}

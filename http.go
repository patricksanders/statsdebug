package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func serve() {
	http.HandleFunc("/all", httpAllHandler)
	http.HandleFunc("/metric/", httpMetricHandler)
	http.HandleFunc("/reset", httpResetHandler)
	srv := &http.Server{
		Addr:         ":8080",
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Info("Stopped serving", "err", err)
		}
	}()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	log.Info("Shutting down")
}

func httpAllHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("handling http request", "path", r.URL.Path)
	json.NewEncoder(w).Encode(summary.getAll())
}

func httpMetricHandler(w http.ResponseWriter, r *http.Request) {
	stat := r.URL.Path[8:]
	log.Info("handling http request", "path", r.URL.Path)
	json.NewEncoder(w).Encode(summary.get(stat))
}

func httpResetHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("handling http request", "path", r.URL.Path)
	summary.reset()
}

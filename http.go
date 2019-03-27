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
	http.HandleFunc("/", httpStatHandler)
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

type CountResponse struct {
	Count int `json:"count"`
}

func httpStatHandler(w http.ResponseWriter, r *http.Request) {
	stat := r.URL.Path[1:]
	log.Info("handling http request", "path", stat)
	response := CountResponse{
		Count: counter.get(stat),
	}
	json.NewEncoder(w).Encode(&response)
}

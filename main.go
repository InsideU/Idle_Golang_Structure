package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/InsideU/go-module/handlers"
)

func main() {

	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	hh := handlers.NewHello(l)
	ph := handlers.NewProducts(l)
	gh := handlers.NewGoodbye(l)
	sm := http.NewServeMux()
	sm.Handle("/", hh)
	sm.Handle("/goodbye", gh)
	sm.Handle("/products", ph)

	s := &http.Server{
		Addr:         ":8080",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill, syscall.SIGTERM)

	sig := <-sigChan
	l.Println("Received graceful termination ", sig)
	//gracefull shutdown in case when we want to upgrade the system we will let the current work
	//of the client get finished and then shutdown the system so that no task is blocked
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc) // wait until all the requst is completed but will not accept any new request

	http.ListenAndServe(":8080", sm)
}

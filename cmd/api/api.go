package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type (
	application struct {
		conf config
		logger *zap.SugaredLogger
	}

	config struct {
		addr string
		env  string
	}
)

func (app *application) mount() http.Handler{
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped
	r.Use(middleware.Timeout(60 * time.Second))


	// use Grop routing feature of chi
	r.Route("/v1", func(r chi.Router){
		r.Get("/health", app.healthCheckHandler)
	})
	return r

}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.conf.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("error cannot run server")
		return err
	}
	return nil

}

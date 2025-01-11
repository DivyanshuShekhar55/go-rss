package main

import (
	"log"
	"net/http"
	"time"

	"github.com/DivyanshuShekhar55/go-rss/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

type application struct {
	conf   config
	logger *zap.SugaredLogger
	store  store.Storage

}

type config struct {
	addr string
	env  string
	db   dbConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped
	r.Use(middleware.Timeout(60 * time.Second))

	// use Group routing feature of chi
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.HealthCheckHandler)

		r.Route("/rss", func(r chi.Router) {
			r.Get("/get", app.GetFeedHandler)
		})

		r.Route("/users", func(r chi.Router){
			
		})
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

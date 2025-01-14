package main

import (
	"log"
	"net/http"
	"time"

	"github.com/DivyanshuShekhar55/go-rss/internal/auth"
	"github.com/DivyanshuShekhar55/go-rss/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

type application struct {
	conf         config
	logger       *zap.SugaredLogger
	store        store.Storage
	autheticator auth.Authenticator
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

// will be used in jwt authentication
type basicConfig struct{
	user string
	pass string
}

// will be used in jwt authentication
type tokenConfig struct{
	secret string
	exp time.Duration
	iss string
}

// will be used in jwt authentication 
type authConfig struct{
	basic basicConfig
	token tokenConfig
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	//following midddlewares will be used by all the routes
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

			// http://localhost:8.../v1/rss/get
			r.Get("/get", app.GetFeedHandler)
		})

		r.Route("/users", func(r chi.Router) {
			// public route activate/{token}
			r.Put("/activate/{token}", app.activateUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				// all the routes like .../users/{userID}/* should be authenticated protected
				r.Use(app.AuthTokenMiddleware)
				r.Get("/", app.getUserHandler)
			})
		})

		// Public routes ...
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)
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

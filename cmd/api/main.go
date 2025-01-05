package main

import (
	"fmt"

	"github.com/DivyanshuShekhar55/go-rss/internal/env"
	"go.uber.org/zap"
)

func main() {
	cfg := config{
		addr: env.GetString("PORT", ":8080"),
		env:  env.GetString("ENV", "dev"),
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	app := &application{
		conf:   cfg,
		logger: logger,
	}

	mux := app.mount()
	fmt.Println("running")
	logger.Fatal(app.run(mux))

}

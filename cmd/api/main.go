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
		db : dbConfig{
			addr:"",
			maxOpenConns: 30,
			maxIdleConns: 30,
			maxIdleTime: "15ms",
		},
	}

	mux := app.mount()
	fmt.Println("running")
	logger.Fatal(app.run(mux))

}

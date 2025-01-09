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
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", ""),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
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

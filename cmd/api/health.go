package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "ok",
		"env":    app.conf.env,
	}

	if err := app.jsonResponse(w, data, http.StatusOK); err != nil {
		return
	}
	fmt.Println(data)
}

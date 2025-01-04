package main

import (
	"fmt"
	"net/http"
)

func(app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request){
	data := map[string]string{
		"status" : "ok",
		"env": app.conf.env,
	}

	// if err:= app.jsonResponse(w, http.StatusOK, data); err != nil {
	// 	return
	// }
	fmt.Println(data)
}
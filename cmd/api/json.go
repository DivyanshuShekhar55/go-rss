// https://github.com/sikozonpc/GopherSocial/blob/dev/cmd/api/json.go FOR MORE THINGS

package main

import (
	"encoding/json"
	"net/http"
)

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {

	maxBytes := 1_048_578 //1mb
	r.Body = (http.MaxBytesReader(w, r.Body, int64(maxBytes)))

	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(data)
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func (app *application) jsonResponse(w http.ResponseWriter, data any, status int) error {
	type envelope struct {
		Data any `json:"data"`
	}

	return writeJSON(w, status, &envelope{Data: data})
}

package main

import "net/http"

type RegisteredUserPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {

}

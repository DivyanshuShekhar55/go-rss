package main

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/DivyanshuShekhar55/go-rss/internal/store"
)

var validate *validator.Validate

type RegisteredUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	validate = validator.New()

	var payload RegisteredUserPayload
	if err := readJSON(w, r, &payload); err !=nil{
		app.badRequestResponse(w, r, err)
		return
	}

	// do some validation checks
	if err:=validate.Struct(payload); err!=nil{
		app.badRequestResponse(w, r, err)
		return
	}

	// user has been validated, all good in struct, can now try to put to db
	user := &store.User{
		Username: payload.Username,
		Email:payload.Email,
		// put the password after hashing it
	}
	
	if err:= user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	


}

func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {

}

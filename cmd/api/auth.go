package main

import (
	"net/http"

	"github.com/DivyanshuShekhar55/go-rss/internal/store"
	"github.com/go-playground/validator/v10"
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
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// do some validation checks
	if err := validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// user has been validated, all good in struct, can now try to put to db
	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
		// put the password after hashing it
	}

	// HASHING THE PASSWORD
	// type User struct {
	// 	ID        int64    `json:"id"`
	// 	Username  string   `json:"username"`
	// 	Email     string   `json:"email"`
	// 	Password  password `json:"-"`
	// 	CreatedAt string   `json:"created_at"`
	// }
	// here is how the user struct looks like and note that the passwords struct has a Set function (refer to users.go file in internals/store)
	// the user var created locally here, will get the password field set as user.Password = ""...
	// but rather we are saying the password will pass through the Set function and then be set ... basically calling function on a struct field

	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	// generate the token, here if you want to have the functionality to activate the user via email sending
	// put functionality of email, sending here if you want it

}

func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {

}

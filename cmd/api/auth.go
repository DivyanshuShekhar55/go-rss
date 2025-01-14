package main

import (
	"net/http"
	"time"

	"github.com/DivyanshuShekhar55/go-rss/internal/store"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
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

	//ctx := r.Context()

	// generate the token, here if you want to have the functionality to activate the user via email sending
	// put functionality of email, sending here if you want it

	// CREATE the user in the db
	// TODO

}

type CreateUserTokenPayload struct {
	// note that the user id in the store has the json of id only not ID
	ID       int64  `json:"id" validate:"required"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

// LOGIN like system where we check the login details with the db details
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {

	// read the json body of payload
	var payload CreateUserTokenPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// validate against the struct whether it is properly "formatted"
	if err := validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// the user payload is fine, can procced to check the user details against the db details
	// ideally it should be getUserByEmail as the user will put his/her mail only in frontend
	user, err := app.store.Users.GetByID(r.Context(), payload.ID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.unauthorisedErrorResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	// user with the given id found
	// proceed to compare password
	if err := user.Password.Compare(payload.Password); err != nil {
		app.unauthorisedErrorResponse(w, r, err)
		return
	}

	// password also correct, move to creating new jwt token
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.conf.auth.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.conf.auth.token.iss,

		// *** HAVE A DOUBT ON THE FOLLOWING PART .iss or .aud ?? ***
		"aud": app.conf.auth.token.iss,
	}

	// jwt has been structured, now generate it
	token, err := app.autheticator.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// token generated, send response
	if err := app.jsonResponse(w, token, http.StatusCreated); err != nil {
		app.internalServerError(w, r, err)
	}

}

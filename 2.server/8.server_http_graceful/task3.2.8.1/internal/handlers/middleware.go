package handlers

import (
	"net/http"

	"github.com/go-chi/jwtauth"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return jwtauth.Verifier(TokenAuth)(jwtauth.Authenticator(next))
}

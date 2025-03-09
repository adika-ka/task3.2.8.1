package handlers

import "github.com/go-chi/jwtauth"

var TokenAuth = jwtauth.New("HS256", []byte("secret"), nil)

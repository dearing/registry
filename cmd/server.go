package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dearing/registry"
)

type server struct {
	*http.ServeMux
	*registry.Registry
}

var secret = []byte("supersecretkey")

func NewServer() *server {
	s := &server{
		ServeMux: http.NewServeMux(),
		Registry: registry.NewRegistry(),
	}

	s.HandleFunc("/register", s.handleRegister)
	s.HandleFunc("/login", s.handleLogin)
	s.HandleFunc("/session", s.handleSession)

	return s
}

func (s *server) handleRegister(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// retrieve the username and password from the request body
	username := r.FormValue("username")
	password := r.FormValue("password")

	// register the username and password
	err := s.Register(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// successful registration
	http.Error(w, fmt.Sprintf("User %s Registered", username), http.StatusOK)
}

func (s *server) handleLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// retrieve the username and password from the request body
	username := r.FormValue("username")
	password := r.FormValue("password")

	err := s.Verify(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// create JWT
	jwt := registry.NewJWT(username, time.Now().Add(time.Hour*24).Unix())
	if err != nil {
		println(err.Error())
		// return error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// encode JWT
	token, err := jwt.Encode(secret)
	if err != nil {
		println(err.Error())
		// return error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create a secure cookie with the JWT token
	cookie := http.Cookie{
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		Name:     "jwt",
		SameSite: http.SameSiteStrictMode,
		Secure:   false,
		Value:    token,
	}

	http.SetCookie(w, &cookie)

	message := fmt.Sprintf("Welcome back %s, your login will expire at %s\n", jwt.Payload.Subject, jwt.Payload.ExpiresAt.HumanReadable())
	w.Write([]byte(message))
}

func (s *server) handleSession(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// retrieve the JWT from the secure cookie
	cookie, err := r.Cookie("jwt")
	if err != nil {
		println(err.Error())
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// decode the JWT token
	jwt, err := registry.Decode(cookie.Value, secret)
	if err != nil {
		println(err.Error())
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	jsonData, err := json.Marshal(jwt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return stored jwt
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

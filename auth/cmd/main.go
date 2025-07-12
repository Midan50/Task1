package main

import (
	"log"
	"net/http"

	"auth/internal/config"
	"auth/internal/handler"
	"auth/internal/middleware"
	"auth/internal/storage"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.MustLoad()
	db := storage.MustConnect(cfg.DBUrl)
	defer db.Close()

	r := chi.NewRouter()
	h := handler.NewHandler(db, cfg)

	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	r.With(middleware.Auth(cfg.JWTSecret)).Get("/logout", h.Logout)

	log.Printf("Server running at %s:%s", cfg.Host, cfg.Port)
	http.ListenAndServe(cfg.Host+":"+cfg.Port, r)
}

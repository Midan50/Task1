package handler

import (
	"auth/internal/config"
	"auth/internal/model"
	"context"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

func NewHandler(db *pgxpool.Pool, cfg *config.Config) *Handler {
	return &Handler{db: db, cfg: cfg}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var u model.User
	render.DecodeJSON(r.Body, &u)
	_, err := h.db.Exec(context.Background(), "INSERT INTO users (email, password) VALUES ($1, $2)", u.Email, u.Password)
	if err != nil {
		http.Error(w, "registration error", http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, map[string]string{"status": "registered"})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var u model.User
	render.DecodeJSON(r.Body, &u)
	row := h.db.QueryRow(context.Background(), "SELECT id, password FROM users WHERE email=$1", u.Email)
	var id int
	var pw string
	if err := row.Scan(&id, &pw); err != nil || pw != u.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	})
	token, err := claims.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, map[string]string{"token": token})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, map[string]string{"status": "logged out"})
}

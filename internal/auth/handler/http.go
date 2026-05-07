package handler

import (
	"EventSpace/internal/auth/service"
	"EventSpace/pkg/dto"
	"encoding/json"
	"log"
	"net/http"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	tokens, err := h.authService.Register(r.Context(), req.Email, req.Password)

	if err != nil {
		log.Printf("Registration error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.setRefreshTokenCookie(w, tokens.RefreshToken)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token": tokens.AccessToken,
		"expires_at":   tokens.ExpiresAt,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	tokens, err := h.authService.Login(r.Context(), req.Email, req.Password)

	if err != nil {
		log.Printf("Registration error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.setRefreshTokenCookie(w, tokens.RefreshToken)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token": tokens.AccessToken,
		"expires_at":   tokens.ExpiresAt,
	})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "missing refresh token", http.StatusUnauthorized)
		return
	}

	tokens, err := h.authService.RefreshAccessToken(r.Context(), cookie.Value)

	if err != nil {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		log.Fatalf("Error: %v", err.Error())
		return
	}

	h.setRefreshTokenCookie(w, tokens.RefreshToken)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token": tokens.AccessToken,
		"expires_at":   tokens.ExpiresAt,
	})

}

func (h *AuthHandler) setRefreshTokenCookie(w http.ResponseWriter, refreshToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/refresh",
		MaxAge:   7 * 24 * 3600,
	})
}

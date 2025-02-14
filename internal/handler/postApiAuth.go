package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"merch-shop/internal/api"
	internalErrors "merch-shop/internal/errors"
	"merch-shop/internal/httputil"
	"merch-shop/pkg/jwt"
	"merch-shop/pkg/logger"
)

func (h *Handler) PostApiAuth(w http.ResponseWriter, r *http.Request) {
	var req api.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetOrCreate(r.Context(), req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, internalErrors.ErrUserCreationFailed):
			httputil.RespondError(w, "Failed to create user", http.StatusInternalServerError)
			return
		case errors.Is(err, internalErrors.ErrInvalidCredentials):
			httputil.RespondError(w, "Invalid credentials", http.StatusUnauthorized)
			return
		default:
			logger.Error("PostApiAuth: ", "err", err)
			httputil.RespondError(w, "Failed to create user", http.StatusInternalServerError)
		}
	}

	token, err := jwt.GenerateToken(user.ID, h.secret)
	if err != nil {
		logger.Error("PostApiAuth jwt.GenerateToken: ", "err", err)
		httputil.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httputil.RespondJSON(w, api.AuthResponse{Token: &token}, http.StatusOK)
}

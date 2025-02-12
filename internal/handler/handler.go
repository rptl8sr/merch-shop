package handler

import (
	"encoding/json"
	"net/http"

	"merch-shop/internal/api"
	"merch-shop/internal/httputil"
	mdlwr "merch-shop/internal/middleware"
	"merch-shop/internal/model"
	"merch-shop/pkg/jwt"
)

var (
	success = map[string]string{"status": "success"}
)

type UserService interface {
	GetOrCreate(username, password string) (model.User, error)
}

type InfoService interface {
	GetUserInfo(userID uint) (api.InfoResponse, error)
	BuyItem(userID uint, item string) error
	SendCoins(userID uint, toUser string, amount int) error
}

type Handler struct {
	secret      string
	userService UserService
	infoService InfoService
}

func New(secret string, userService UserService, infoService InfoService) *Handler {
	return &Handler{
		secret:      secret,
		userService: userService,
		infoService: infoService,
	}
}

func (h *Handler) PostApiAuth(w http.ResponseWriter, r *http.Request) {
	var req api.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetOrCreate(req.Username, req.Password)
	if err != nil {
		httputil.RespondError(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	token, err := jwt.GenerateToken(user.ID, h.secret)
	if err != nil {
		httputil.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httputil.RespondJSON(w, api.AuthResponse{Token: &token}, http.StatusOK)
}

func (h *Handler) GetApiInfo(w http.ResponseWriter, r *http.Request) {
	userID, err := mdlwr.GetUserIDFromContext(r.Context())
	if err != nil {
		httputil.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	info, err := h.infoService.GetUserInfo(userID)
	if err != nil {
		httputil.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httputil.RespondJSON(w, info, http.StatusOK)
}

func (h *Handler) GetApiBuyItem(w http.ResponseWriter, r *http.Request, item string) {
	userID, err := mdlwr.GetUserIDFromContext(r.Context())
	if err != nil {
		httputil.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.infoService.BuyItem(userID, item)
	if err != nil {
		httputil.RespondError(w, "Failed to buy item: "+err.Error(), http.StatusBadRequest)
		return
	}

	httputil.RespondJSON(w, success, http.StatusOK)
}

func (h *Handler) PostApiSendCoin(w http.ResponseWriter, r *http.Request) {
	var req api.SendCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := mdlwr.GetUserIDFromContext(r.Context())
	if err != nil {
		httputil.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.infoService.SendCoins(userID, req.ToUser, req.Amount)
	if err != nil {
		httputil.RespondError(w, "Failed to send coins: "+err.Error(), http.StatusBadRequest)
		return
	}

	httputil.RespondJSON(w, success, http.StatusOK)
}

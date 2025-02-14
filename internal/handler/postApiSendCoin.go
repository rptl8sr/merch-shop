package handler

import (
	"encoding/json"
	"net/http"

	"merch-shop/internal/api"
	"merch-shop/internal/httputil"
	mdlwr "merch-shop/internal/middleware"
)

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

	err = h.transactionService.SendCoins(r.Context(), userID, req.ToUser, req.Amount)
	if err != nil {
		httputil.RespondError(w, "Failed to send coins: "+err.Error(), http.StatusBadRequest)
		return
	}

	httputil.RespondJSON(w, success, http.StatusOK)
}

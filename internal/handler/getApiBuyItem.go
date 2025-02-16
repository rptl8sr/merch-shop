package handler

import (
	"net/http"

	"merch-shop/internal/httputil"
	mdlwr "merch-shop/internal/middleware"
)

func (h *Handler) GetApiBuyItem(w http.ResponseWriter, r *http.Request, item string) {
	userID, err := mdlwr.GetUserIDFromContext(r.Context())
	if err != nil {
		httputil.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.purchaseService.BuyItem(r.Context(), userID, item, &defaultQuantity)
	if err != nil {
		httputil.RespondError(w, "Failed to buy item: "+err.Error(), http.StatusBadRequest)
		return
	}

	httputil.RespondJSON(w, success, http.StatusOK)
}

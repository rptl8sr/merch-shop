package handler

import (
	"net/http"

	"merch-shop/internal/httputil"
	mdlwr "merch-shop/internal/middleware"
)

func (h *Handler) GetApiInfo(w http.ResponseWriter, r *http.Request) {
	userID, err := mdlwr.GetUserIDFromContext(r.Context())
	if err != nil {
		httputil.RespondError(w, "Failed to get user: "+err.Error(), http.StatusBadRequest)
		return
	}

	info, err := h.infoService.GetUserInfo(r.Context(), userID)
	if err != nil {
		httputil.RespondError(w, "Failed to get info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	httputil.RespondJSON(w, info, http.StatusOK)
}

package handler

import (
	"encoding/json"
	"github.com/kstsm/wb-event-booker/internal/dto"
	"net/http"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := req.ValidateUser(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.CreateUser(r.Context(), &req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, dto.CreateUserResponse{
		User:    user,
		Message: "user created successfully",
	})
}

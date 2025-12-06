package handler

import (
	"encoding/json"
	"errors"
	"github.com/kstsm/wb-event-booker/internal/apperrors"
	"github.com/kstsm/wb-event-booker/internal/dto"
	"net/http"
)

func (h *Handler) createEventHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEventRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := req.ValidateEvent(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	event, err := h.service.CreateEvent(r.Context(), &req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, dto.CreateEventResponse{
		Event:   event,
		Message: "event created successfully",
	})
}

func (h *Handler) getEventByIDHandler(w http.ResponseWriter, r *http.Request) {
	eventID, err := parseUUIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	event, err := h.service.GetEventByID(r.Context(), eventID)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.EventNotFound):
			respondError(w, http.StatusNotFound, "event not found")
		default:
			respondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondJSON(w, http.StatusOK, dto.GetEventResponse{
		Event: event,
	})
}

func (h *Handler) listEventsHandler(w http.ResponseWriter, r *http.Request) {
	events, err := h.service.ListEvents(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, dto.ListEventsResponse{
		Events: events,
	})
}

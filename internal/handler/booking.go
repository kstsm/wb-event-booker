package handler

import (
	"encoding/json"
	"errors"
	"github.com/kstsm/wb-event-booker/internal/apperrors"
	"github.com/kstsm/wb-event-booker/internal/dto"
	"net/http"
)

func (h *Handler) bookEventHandler(w http.ResponseWriter, r *http.Request) {
	eventID, err := parseUUIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req dto.BookEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Email == "" {
		respondError(w, http.StatusBadRequest, "email is required")
		return
	}

	booking, err := h.service.BookEvent(r.Context(), eventID, &req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.EventNotFound):
			respondError(w, http.StatusNotFound, "event not found")
		case errors.Is(err, apperrors.UserNotFound):
			respondError(w, http.StatusNotFound, "user not found")
		case errors.Is(err, apperrors.NoAvailableSeats):
			respondError(w, http.StatusConflict, "no available seats")
		case errors.Is(err, apperrors.UserAlreadyBookedThisEvent):
			respondError(w, http.StatusConflict, "user already has a booking for this event")
		case errors.Is(err, apperrors.EventExpired):
			respondError(w, http.StatusBadRequest, "event has expired")
		default:
			respondError(w, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	respondJSON(w, http.StatusCreated, dto.BookEventResponse{
		BookingID: booking.BookingID,
		Deadline:  booking.Deadline,
		Message:   "booking created successfully",
	})
}

func (h *Handler) ConfirmBookingHandler(w http.ResponseWriter, r *http.Request) {
	eventID, err := parseUUIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req dto.ConfirmBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.ConfirmBooking(r.Context(), eventID, &req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.BookingNotFound):
			respondError(w, http.StatusNotFound, "booking not found")
		case errors.Is(err, apperrors.BookingNotReserved):
			respondError(w, http.StatusBadRequest, "booking is not in reserved status")
		case errors.Is(err, apperrors.BookingDeadlinePassed):
			respondError(w, http.StatusBadRequest, "booking deadline has passed")
		case errors.Is(err, apperrors.EventDoesNotRequirePayment):
			respondError(w, http.StatusBadRequest, "event does not require payment confirmation")
		case errors.Is(err, apperrors.EventExpired):
			respondError(w, http.StatusBadRequest, "event has expired")
		default:
			respondError(w, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	respondJSON(w, http.StatusOK, dto.ConfirmBookingResponse{
		Message: "confirmed successfully",
	})
}

func (h *Handler) listBookingsByEventHandler(w http.ResponseWriter, r *http.Request) {
	eventID, err := parseUUIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	bookings, err := h.service.ListBookingsByEventID(r.Context(), eventID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, dto.ListBookingsResponse{
		Bookings: bookings,
	})
}

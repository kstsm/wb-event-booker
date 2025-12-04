package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/kstsm/wb-event-booker/internal/service"
	"net/http"
)

type HandlerI interface {
	NewRouter() http.Handler
}

type Handler struct {
	service service.ServiceI
}

func NewHandler(service service.ServiceI) HandlerI {
	return &Handler{
		service: service,
	}
}

func (h *Handler) NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", h.serveHTML("index.html"))
	r.Get("/register", h.serveHTML("register.html"))
	r.Get("/admin", h.serveHTML("admin.html"))
	r.Get("/event", h.serveHTML("event.html"))

	r.Route("/api", func(r chi.Router) {
		r.Post("/events", h.createEventHandler)
		r.Post("/events/{id}/book", h.bookEventHandler)
		r.Post("/events/{id}/confirm", h.ConfirmBooking)
		r.Get("/events/{id}", h.getEventByIDHandler)

		r.Get("/events", h.ListEvents)
		r.Get("/events/{id}/bookings", h.ListBookingsByEvent)
		r.Post("/users", h.CreateUser)
	})

	return r
}

func (h *Handler) serveHTML(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/"+filename)
	}
}

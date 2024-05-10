package appointments

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/razdacoder/mcwale-api/utils"
)

type Handler struct {
	store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{
		store: store,
	}
}

func appointmentsRoues(handler *Handler) chi.Router {
	router := chi.NewRouter()

	router.Route("/", func(router chi.Router) {
		router.Get("/", handler.HandleGetAllAppointments)
		router.Post("/", handler.HandleCrateAppointment)
	})

	router.Route("/{id}", func(route chi.Router) {
		route.Get("/", handler.HandleSingleAppointment)
	})

	return router
}

func (handler *Handler) RegisterRoutes(router chi.Router) {
	router.Mount("/appointments", appointmentsRoues(handler))
}

func (handler *Handler) HandleGetAllAppointments(writer http.ResponseWriter, request *http.Request) {
	appointments, err := handler.store.GetAllAppointments()
	if err != nil {
		utils.WriteError(writer, http.StatusInternalServerError,
			fmt.Errorf("internal server error"))
		return
	}

	utils.WriteJSON(writer, http.StatusOK, appointments)
}

func (handler *Handler) HandleCrateAppointment(writer http.ResponseWriter, request *http.Request) {
	var payload CreateAppointmentPayload
	if err := utils.ParseJSON(request, &payload); err != nil {
		utils.WriteError(writer, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(writer, http.StatusUnprocessableEntity, errors)
		return
	}
	err := handler.store.CreateNewAppointments(payload)
	if err != nil {
		utils.WriteError(writer, http.StatusInternalServerError,
			fmt.Errorf("internal server error"))
		return
	}

	utils.WriteJSON(writer, http.StatusOK, map[string]string{"message": "Appointment Booked"})
}

func (handler *Handler) HandleSingleAppointment(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	if id == "" {
		utils.WriteError(writer, http.StatusBadRequest,
			fmt.Errorf("no id found"))
		return
	}
	app, err := handler.store.GetSingleAppointment(id)
	if err != nil {
		utils.WriteError(writer, http.StatusInternalServerError,
			err)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, app)
}

package orders

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"github.com/razdacoder/mcwale-api/services/auth"
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

func orderRoutes(handler *Handler) chi.Router {
	router := chi.NewRouter()

	router.Route("/", func(router chi.Router) {
		router.Post("/", handler.handleCreateOrder)
		router.Get("/", handler.handleGetOrders)
	})

	router.Route("/{id}", func(router chi.Router) {
		router.Get("/", handler.handleGetOrder)
		router.With(auth.IsLoggedIn, auth.IsAdmin).Patch("/", handler.handleUpdateOrderStatus)
		router.With(auth.IsLoggedIn, auth.IsAdmin).Delete("/", handler.handleDeleteOrder)
	})

	return router
}

func (handler *Handler) RegisterRoutes(router chi.Router) {
	router.Mount("/orders", orderRoutes(handler))
}

func (handler *Handler) handleCreateOrder(writer http.ResponseWriter, request *http.Request) {
	var payload CreateOrderPayload
	if err := utils.ParseJSON(request, &payload); err != nil {
		utils.WriteError(writer, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	orderNumber, err := handler.store.CreateOrder(payload)

	if err != nil {
		utils.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(writer, http.StatusCreated, map[string]string{"order_number": orderNumber})
}

func (handler *Handler) handleGetOrders(writer http.ResponseWriter, request *http.Request) {
	orders, err := handler.store.GetOrders()
	if err != nil {
		utils.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, orders)
}

func (handler *Handler) handleGetOrder(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	if id == "" {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("no order id"))
		return
	}
	order, err := handler.store.GetOrderByID(id)
	if err != nil {
		utils.WriteError(writer, http.StatusNotFound, err)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, order)
}

func (handler *Handler) handleUpdateOrderStatus(writer http.ResponseWriter, request *http.Request) {

	id := chi.URLParam(request, "id")
	if id == "" {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("no order id"))
		return
	}

	var payload UpdateOrderPayload
	if err := utils.ParseJSON(request, &payload); err != nil {
		utils.WriteError(writer, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	err := handler.store.UpdateOrderStatus(id, payload.Status)
	if err != nil {
		utils.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, map[string]string{"message": "Order Status updated"})

}

func (handler *Handler) handleDeleteOrder(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	if id == "" {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("no order id"))
		return
	}
	err := handler.store.DeleteOrder(id)
	if err != nil {
		utils.WriteError(writer, http.StatusNotFound, err)
		return
	}
	utils.WriteJSON(writer, http.StatusNoContent, map[string]string{"message": "Order Deleted"})
}

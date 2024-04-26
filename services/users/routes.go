package users

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"github.com/razdacoder/mcwale-api/models"
	"github.com/razdacoder/mcwale-api/services/auth"
	"github.com/razdacoder/mcwale-api/utils"
)

type Handler struct {
	store UserStore
}

func NewHandler(store UserStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (handler *Handler) RegisterRoutes(router chi.Router) {
	router.Post("/login", handler.handleLogin)
	router.Post("/register", handler.handleRegister)
	router.Get("/users", handler.handleGetAllUsers)
	router.Get("/users/{id}", handler.handleGetSingleUser)
	router.Put("/users/{id}", handler.handleUpdateUser)
	router.Delete("/users/{id}", handler.handleUserDelete)
}

func (handler *Handler) handleLogin(writer http.ResponseWriter, request *http.Request) {
	var payload LoginUserPayload
	if err := utils.ParseJSON(request, &payload); err != nil {
		utils.WriteError(writer, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	user, err := handler.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(writer, http.StatusUnauthorized, fmt.Errorf("invalid email or password"))
		return
	}

	if !auth.VerifyPassword(user.Password, []byte(payload.Password)) {
		utils.WriteError(writer, http.StatusUnauthorized, fmt.Errorf("invalid email or password"))
		return
	}

	token, err := auth.CreateJWT([]byte(os.Getenv("JWT_SECRET")), user.ID)
	if err != nil {
		utils.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, map[string]string{"token": token})
}

func (handler *Handler) handleRegister(writer http.ResponseWriter, request *http.Request) {
	var payload RegisterUserPayload
	if err := utils.ParseJSON(request, &payload); err != nil {
		utils.WriteError(writer, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	userExists, err := handler.store.UserExists(payload.Email)
	if userExists || err != nil {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("user with email address %s already exists", payload.Email))
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	err = handler.store.CreateUser(RegisterUserPayload{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})

	if err != nil {
		utils.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(writer, http.StatusCreated, nil)
}

func (handler *Handler) handleGetAllUsers(writer http.ResponseWriter, request *http.Request) {
	users, err := handler.store.GetAllUsers()
	if err != nil {
		utils.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, users)
}

func (handler *Handler) handleGetSingleUser(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	user, err := handler.store.GetUserByID(id)
	if err != nil {
		utils.WriteError(writer, http.StatusNotFound, err)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, user)
}

func (handler *Handler) handleUpdateUser(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	user, err := handler.store.GetUserByID(id)
	if err != nil {
		utils.WriteError(writer, http.StatusNotFound, err)
		return
	}
	var patch models.User
	if err := utils.ParseJSON(request, &patch); err != nil {
		utils.WriteError(writer, http.StatusBadRequest, err)
		return
	}
	patch.Password = user.Password

	err = handler.store.UpdateUserInfo(&patch)
	if err != nil {
		utils.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, patch)
}

func (handler *Handler) handleUserDelete(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	user, err := handler.store.GetUserByID(id)
	if err != nil {
		utils.WriteError(writer, http.StatusNotFound, err)
		return
	}
	err = handler.store.DeleteUser(user.ID)
	if err != nil {
		utils.WriteError(writer, http.StatusBadRequest, err)
	}

	utils.WriteJSON(writer, http.StatusOK, nil)
}

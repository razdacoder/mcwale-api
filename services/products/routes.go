package products

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"github.com/razdacoder/mcwale-api/models"
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

func categoriesRouter(handler *Handler) chi.Router {
	router := chi.NewRouter()

	router.Route("/", func(router chi.Router) {
		router.Get("/", handler.handleGetAllCategories)
		router.Post("/", handler.handleCreateCategory)
	})

	router.Route("/{slug}", func(router chi.Router) {
		router.Get("/", handler.handleGetSingleCategory)
		router.With(auth.IsLoggedIn, auth.IsAdmin).Patch("/", handler.handleUpdateCategory)
		router.With(auth.IsLoggedIn, auth.IsAdmin).Delete("/", handler.handleDeleteCategory)
	})

	return router
}

func productsRouter(handler *Handler) chi.Router {
	router := chi.NewRouter()

	router.Get("/category/{slug}", handler.handleGetProductsByCategory)
	router.Get("/recent", handler.handleGetRecentProducts)
	router.Get("/featured", handler.handleGetFeaturedProducts)

	router.Route("/", func(router chi.Router) {
		router.Get("/", handler.handleGetAllProducts)

		router.With(auth.IsLoggedIn, auth.IsAdmin).Post("/", handler.handleCreateProduct)
	})

	router.Route("/{slug}", func(router chi.Router) {
		router.Get("/", handler.handleGetSingleProduct)
		router.With(auth.IsLoggedIn, auth.IsAdmin).Patch("/", handler.handleUpdateProduct)
		router.With(auth.IsLoggedIn, auth.IsAdmin).Delete("/", handler.handleDeleteProduct)
	})

	return router
}

func (handler *Handler) RegisterRoutes(router chi.Router) {
	router.Mount("/categories", categoriesRouter(handler))
	router.Mount("/products", productsRouter(handler))
}

func (handler *Handler) handleGetAllCategories(writer http.ResponseWriter, request *http.Request) {
	categories, err := handler.store.GetAllCategories()
	if err != nil {
		utils.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, categories)
}

func (handler *Handler) handleCreateCategory(writer http.ResponseWriter, request *http.Request) {
	var payload CreateCategoryPayload
	if err := utils.ParseJSON(request, &payload); err != nil {
		utils.WriteError(writer, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	if err := handler.store.CreateCategory(payload); err != nil {
		utils.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(writer, http.StatusCreated, map[string]string{"message": "Category Created"})
}

func (handler *Handler) handleGetSingleCategory(writer http.ResponseWriter, request *http.Request) {
	slug := chi.URLParam(request, "slug")
	if slug == "" {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("no slug found"))
		return
	}
	category, err := handler.store.GetSingleCategory(slug)
	if err != nil {
		utils.WriteError(writer, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(writer, http.StatusCreated, category)
}

func (handler *Handler) handleUpdateCategory(writer http.ResponseWriter, request *http.Request) {
	slug := chi.URLParam(request, "slug")
	if slug == "" {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("no slug found"))
		return
	}

	var patch models.Category
	if err := utils.ParseJSON(request, &patch); err != nil {
		utils.WriteError(writer, http.StatusBadRequest, err)
		return
	}

	err := handler.store.UpdateCategory(slug, &patch)
	if err != nil {
		utils.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, map[string]string{"message": "Updated Successfully"})
}

func (handler *Handler) handleDeleteCategory(writer http.ResponseWriter, request *http.Request) {
	slug := chi.URLParam(request, "slug")
	if slug == "" {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("no slug found"))
		return
	}
	err := handler.store.DeleteCategory(slug)
	if err != nil {
		utils.WriteError(writer, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(writer, http.StatusNoContent, map[string]string{"message": "Deleted Successfully"})
}

func (handler *Handler) handleCreateProduct(writer http.ResponseWriter, request *http.Request) {
	var payload CreateProductPayload
	err := utils.ParseJSON(request, &payload)
	if err != nil {
		utils.WriteError(writer, http.StatusBadRequest, err)
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	if err := handler.store.CreateProduct(payload); err != nil {
		utils.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(writer, http.StatusCreated, map[string]string{"message": "Product Created"})
}

func (handler *Handler) handleGetAllProducts(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	style := query.Get("style")
	minPrice, _ := strconv.ParseFloat(query.Get("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(query.Get("max_price"), 64)
	category_slug := query.Get("category")
	page := utils.ParseStringToInt(query.Get("page"), 0)
	perPage := utils.ParseStringToInt(os.Getenv("PER_PAGE"), 10)
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * perPage
	products, err := handler.store.GetAllProducts(style, category_slug, minPrice, maxPrice, offset)
	if err != nil {
		utils.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, map[string]any{"results": len(products), "page": page, "data": products})
}

func (handler *Handler) handleGetSingleProduct(writer http.ResponseWriter, request *http.Request) {
	slug := chi.URLParam(request, "slug")
	if slug == "" {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("no slug found"))
		return
	}
	product, err := handler.store.GetSingleProduct(slug)
	if err != nil {
		utils.WriteError(writer, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(writer, http.StatusOK, product)
}

func (handler *Handler) handleUpdateProduct(writer http.ResponseWriter, request *http.Request) {
	slug := chi.URLParam(request, "slug")
	if slug == "" {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("no slug found"))
		return
	}

	var patch models.Product
	if err := utils.ParseJSON(request, &patch); err != nil {
		utils.WriteError(writer, http.StatusBadRequest, err)
		return
	}

	err := handler.store.UpdateProduct(slug, &patch)
	if err != nil {
		utils.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, map[string]string{"message": "Updated Successfully"})
}

func (handler *Handler) handleDeleteProduct(writer http.ResponseWriter, request *http.Request) {
	slug := chi.URLParam(request, "slug")
	if slug == "" {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("no slug found"))
		return
	}
	err := handler.store.DeleteProduct(slug)
	if err != nil {
		utils.WriteError(writer, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(writer, http.StatusNoContent, map[string]string{"message": "Deleted Successfully"})
}

func (handler *Handler) handleGetProductsByCategory(writer http.ResponseWriter, request *http.Request) {
	slug := chi.URLParam(request, "slug")
	if slug == "" {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("no slug found"))
		return
	}
	query := request.URL.Query()
	style := query.Get("style")
	minPrice, _ := strconv.ParseFloat(query.Get("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(query.Get("max_price"), 64)
	page := utils.ParseStringToInt(query.Get("page"), 0)
	perPage := utils.ParseStringToInt(os.Getenv("PER_PAGE"), 10)
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * perPage
	products, err := handler.store.GetProductsByCategory(slug, style, minPrice, maxPrice, offset)
	if err != nil {
		utils.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, map[string]any{"results": len(products), "page": page, "data": products})
}

func (handler *Handler) handleGetRecentProducts(writer http.ResponseWriter, request *http.Request) {
	products, err := handler.store.GetRecentProducts()
	if err != nil {
		utils.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, products)
}

func (handler *Handler) handleGetFeaturedProducts(writer http.ResponseWriter, request *http.Request) {
	products, err := handler.store.GetFeaturedProducts()
	if err != nil {
		utils.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, products)
}

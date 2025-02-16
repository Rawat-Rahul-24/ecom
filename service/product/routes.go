package product

import (
	"ecom/types"
	"ecom/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.handleGetProduct).Methods(http.MethodGet)
	router.HandleFunc("/product", h.handleCreateProduct).Methods(http.MethodPost)
	router.HandleFunc("/product/{id:[0-9]+}", h.handleGetProductById).Methods(http.MethodGet)
}

func (h *Handler) handleGetProduct(w http.ResponseWriter, r *http.Request) {

	products, err := h.store.GetProducts()

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, products)
}

func (h * Handler) handleCreateProduct(w http.ResponseWriter, r * http.Request) {
	var payload types.CreateProductPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	//validate payload

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	//create product
	err := h.store.CreateProduct(types.Product{
		Name: payload.Name,
		Description: payload.Description,
		Quantity: payload.Quantity,
		Image: payload.Image,
		Price: payload.Price,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	utils.WriteJson(w, http.StatusCreated, nil)
}

func (h *Handler) handleGetProductById(w http.ResponseWriter, r *http.Request) {
	pathVars := mux.Vars(r)
	id, _ := strconv.Atoi(pathVars["id"])

	p, err := h.store.GetProductById(id)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	utils.WriteJson(w, http.StatusOK, p)
}
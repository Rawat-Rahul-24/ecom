package cart

import (
	"ecom/service/auth"
	"ecom/types"
	"ecom/utils"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store        types.OrderStore
	productStore types.ProductStore
	userStore    types.UserStore
}

func NewHandler(store types.OrderStore, productStore types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, productStore: productStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(h.handleCheckout, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {

	userId := auth.GetUserIDFromContext(r.Context())
	var cart types.CartCheckoutPayload

	if err := utils.ParseJson(r, &cart); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(cart); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	//get products
	productIds, err := getCartItems(cart.Items)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ps, err := h.productStore.GetProductByIds(productIds)

	orderId, totalPrice, err := h.createOrder(ps, cart.Items, userId)

	utils.WriteJson(w, http.StatusCreated, map[string]any{
		"total_price": totalPrice,
		"orderId":     orderId,
	})
}

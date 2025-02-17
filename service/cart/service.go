package cart

import (
	"ecom/types"
	"fmt"
)

func getCartItems(items []types.CartItem) ([]int, error) {
	productIds := make([]int, len(items))

	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for the product %d", item.ProductID)
		}

		productIds[i] = item.ProductID
	}

	return productIds, nil
}

func (h *Handler) createOrder(ps []types.Product, items []types.CartItem, userId int) (int, float64, error) {

	productMap := make(map[int]types.Product)

	for _, product := range ps {
		productMap[product.ID] = product
	}

	// check if all products are all in stock

	if err := checkIfCartIsInStock(items, productMap); err != nil {
		return 0, 0.0, nil
	}

	//calculate the total price
	totalPrice := calculateTotalPrice(items, productMap)
	
	//update the inventory
	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity
		h.productStore.UpdateProduct(product)
	}

	//create the order and create orderItems
	orderId, err := h.store.CreateOrder(types.Order{
		UserID: userId,
		Total: totalPrice,
		Status: "pending",
		Address: "some address",

	})

	if err != nil {
		return 0, 0, nil
	}

	for _, item := range items {
		h.store.CreateOrderItem(types.OrderItem{
			OrderID: orderId,
			ProductId: item.ProductID,
			Quantity: item.Quantity,
			Price: productMap[item.ProductID].Price,
		})
	}


	return orderId, totalPrice, nil
}

func checkIfCartIsInStock(cartItems []types.CartItem, productMap map[int]types.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range cartItems {
		product, ok := productMap[item.ProductID]

		if !ok {
			return fmt.Errorf("product %d is not available in store, please refresh your cart", item.ProductID)
		}

		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %s is not available in the quantity requested", product.Name)
		}
	}

	return nil
}

func calculateTotalPrice(items []types.CartItem, productMap map[int]types.Product) float64 {
	var total float64

	for _, item := range items {
		product := productMap[item.ProductID]

		total += product.Price * float64(item.Quantity)
	}

	return total
}
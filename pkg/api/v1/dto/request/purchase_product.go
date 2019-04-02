package request

type Product struct {
	ProductID uint `json:"productID"`
	Quantity  uint `json:"quantity"`
}

type PurchaseProduct []Product

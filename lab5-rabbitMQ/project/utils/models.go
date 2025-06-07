package utils

type Order struct {
	Hiker   string `json:"hiker"`
	Product string `json:"product"`
}

const ExchangeName = "orders_exchange"

package main

type Order struct {
	Id          int `json:"id"`
	CustomerID  int `json:"customer_id"`
	TotalAmount int `json:"total_amount"`
}

type OrderItem struct {
	Id          int    `json:"id"`
	OrderID     int    `json:"order_id"`
	ProductName string `json:"product_name"`
	ProductID   int    `json:"product_id"`
	Count       int    `json:"count"`
	TotalPrice  int    `json:"total_price"`
	Status      bool   `json:"status"`
}

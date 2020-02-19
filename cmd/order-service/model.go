package main

import (
	"github.com/jinzhu/gorm"
)

type Order struct {
	gorm.Model
	UserID     uint    `json:"user_id"`
	GrandTotal float64 `json:"grand_total"`
}

type Product struct {
	gorm.Model
	Name  string  `json:"name"`
	Qty   float64 `json:"qty"`
	Price float64 `json:"price"`
}

type OrderProducts struct {
	OrderID     uint    `json:"order_id"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Qty         float64 `json:"qty"`
	SubTotal    float64 `json:"sub_total"`
}

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type OrderDetails struct {
	OrderInfo Order           `json:"order_info"`
	Items     []OrderProducts `json:"items"`
}

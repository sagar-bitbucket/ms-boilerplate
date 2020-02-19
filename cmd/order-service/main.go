package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gitlab.com/scalent/ms-boilerplate/sd"
)

var db *gorm.DB
var err error

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/healthcheck", healthcheck).Methods("GET")
	router.HandleFunc("/create", createOrder).Methods("POST")
	router.HandleFunc("/getOrderByID/{id}", getOrderByID).Methods("GET")
	router.HandleFunc("/getOrdersByUserID/{id}", getOrdersByUserID).Methods("GET")
	db, err = gorm.Open("mysql", "root:password@tcp(192.168.1.31:3306)/ecommerce?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		panic(err)
	}

	defer db.Close()
	portNumber := 8080
	sdConfig, _ := sd.DefaultConfig()
	// fmt.Println(sdConfig)
	sdConfig.ServiceID = "order-service-" + strconv.Itoa(portNumber)
	sdConfig.ServiceName = "order-service"
	sdConfig.ServicePort = portNumber
	sdConfig.Tags = []string{"urlprefix-/getOrderByID", "urlprefix-/getOrdersByUserID", "urlprefix-/order/create strip=/order"}

	err := sdConfig.Register()
	if err != nil {
		panic(err)
	}

	fmt.Println("Server is running:8080")
	fmt.Println(http.ListenAndServe(":8080", router))

}

//Healthcheck
func healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Ok")
}

func createOrder(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("create order API called")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Some error occured while reading a input : ", err)
	}

	type OrderDetails struct {
		UserID        uint            `json:"user_id"`
		OrderProducts []OrderProducts `json:"order_products"`
	}

	// var orderProducts []OrderProducts
	var orderDetails OrderDetails
	err = json.Unmarshal(body, &orderDetails)
	if err != nil {
		fmt.Println("Invalid input : ", err)
	}

	var order Order
	var grandTotal float64

	for index, oProduct := range orderDetails.OrderProducts {
		var product Product
		db.Where("id = ?", oProduct.ProductID).Find(&product)
		orderDetails.OrderProducts[index].ProductName = product.Name
		orderDetails.OrderProducts[index].Price = product.Price
		orderDetails.OrderProducts[index].SubTotal = oProduct.Qty * product.Price
		grandTotal += orderDetails.OrderProducts[index].SubTotal
	}
	order.UserID = orderDetails.UserID
	order.GrandTotal = grandTotal
	db.Create(&order)

	for _, oProduct := range orderDetails.OrderProducts {
		oProduct.OrderID = order.ID
		db.Create(oProduct)
		var product Product
		product.ID = oProduct.ProductID
		db.Model(&product).Update("qty", gorm.Expr("qty-?", 1))
	}

	fmt.Fprintf(rw, "New order created")
}

func getOrderByID(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("getOrderByID API called")
	params := mux.Vars(req)
	orderID, err := params["id"]
	if err == false {
		fmt.Fprintf(rw, "ORDER Id missing")
		return
	}

	var order Order
	var orderDetails OrderDetails
	var orderProducts []OrderProducts

	if err := db.Where("id = ?", orderID).Find(&order).Error; gorm.IsRecordNotFoundError(err) {
		fmt.Fprintf(rw, "Order with ID:? NOT FOUND", orderID)
	}

	if err := db.Where("order_id = ?", orderID).Find(&orderProducts).Error; gorm.IsRecordNotFoundError(err) {
		fmt.Fprintf(rw, "Order with ID:? NOT FOUND", orderID)
	}

	orderDetails.OrderInfo = order
	orderDetails.Items = orderProducts
	output, _ := json.Marshal(orderDetails)
	fmt.Fprintf(rw, string(output))
}

func getOrdersByUserID(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("getOrdersByUserID API called")
	params := mux.Vars(req)
	userID, err := params["id"]
	if err == false {
		fmt.Fprintf(rw, "userID missing")
		return
	}

	var orders []Order
	// var orderDetails OrderDetails

	if err := db.Where("user_id = ?", userID).Find(&orders).Error; gorm.IsRecordNotFoundError(err) {
		fmt.Fprintf(rw, "Orders with UserID:? NOT FOUND", userID)
	}

	output, _ := json.Marshal(orders)
	fmt.Fprintf(rw, string(output))
	// fmt.Println(orders)
	// fmt.Fprintf(rw, orders)
}

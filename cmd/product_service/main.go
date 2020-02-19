package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"gitlab.com/scalent/ms-boilerplate/cmd/responses"
	"gitlab.com/scalent/ms-boilerplate/sd"
)

//Product **
type Product struct {
	gorm.Model
	Name  string  `json:"name"`
	Qty   uint64  `json:"qty"`
	Price float64 `json:"price"`
}

//UpdateStockReq **
type UpdateStockReq struct {
	AddStock    uint64 `json:"add_stock"`
	RemoveStock uint64 `json:"remove_stock"`
	ProductID   uint64 `json:"id"`
}

//UpdateStockResp **
type UpdateStockResp struct {
	Message string `json:"message"`
}

var JWTSecretKey = "secrect"

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/healthcheck", healthcheck).Methods("GET")
	r.HandleFunc("/products/{id}", GetByIDHandler).Methods("GET")
	r.HandleFunc("/products/stock", UpdateStockHandler).Methods("POST")
	r.HandleFunc("/products/", CreateHandler).Methods("POST")
	r.HandleFunc("/products/", ProductListHandler).Methods("GET")
	//r.Use(middleware.ValidateMiddleware)
	portNumber := 8080
	sdConfig, _ := sd.DefaultConfig()
	sdConfig.ServiceID = "product-service-" + strconv.Itoa(portNumber)
	sdConfig.ServiceName = "product-service"
	sdConfig.ServicePort = portNumber
	sdConfig.Tags = []string{"/product", "/product-service"}
	err := sdConfig.Register()
	fmt.Println(err)
	fmt.Println(http.ListenAndServe(":"+strconv.Itoa(portNumber), r))

}

//ValidateMiddleware to authenticate jwt
func ValidateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte(JWTSecretKey), nil
				})
				if error != nil {
					responses.WriteErrorResponse(w, http.StatusBadRequest, "Invalid token")
					return
				}
				if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					w.Header().Set("Access-Control-Allow-Origin", "*")
					w.Header().Set("Access-Control-Allow-Credentials", "true")
					w.Header().Set("Access-Control-Allow-Methods", "POST, PATCH, GET, OPTIONS, PUT, DELETE")
					w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,x-requested-with, XMLHttpRequest, Access-Control-Allow-Methods")
					next.ServeHTTP(w, r)

				}
			}
		} else {
			responses.WriteErrorResponse(w, http.StatusBadRequest, "An Authorization Header is Required!")
		}
	})
}

//Healthcheck
func healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

//CreateHandler creates products
func CreateHandler(w http.ResponseWriter, r *http.Request) {

	product := Product{}
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	resp, err := CreateUser(r.Context(), product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	endpointResp, err := json.Marshal(resp)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(endpointResp)

}

//GetByIDHandler **
func GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, ok := params["id"]
	if ok == false {
		w.Write([]byte("ID NOT FOUND"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resp, err := GetProductByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	endpointResp, err := json.Marshal(resp)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(endpointResp)

}

func UpdateStockHandler(w http.ResponseWriter, r *http.Request) {

	req := UpdateStockReq{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	resp, err := UpdateProductStock(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	endpointResp, err := json.Marshal(resp)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(endpointResp)

}

func ProductListHandler(w http.ResponseWriter, r *http.Request) {

	resp, err := GetProductList(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	endpointResp, err := json.Marshal(resp)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(endpointResp)

}

//NewClientConnection returns mysql collection
func NewClientConnection() *gorm.DB {

	client, err := gorm.Open("mysql", "rahul:password@tcp(127.0.0.1:3306)/ecommerce?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		fmt.Println("Error in Create client connection", err)
		panic("Error In Create Client Connection")
	}

	return client

}

//CreateUser 88
func CreateUser(ctx context.Context, product Product) (*Product, error) {

	dbConn := NewClientConnection()
	defer dbConn.Close()

	createOn := time.Now().In(time.UTC)
	// record create Time
	product.CreatedAt = createOn
	fmt.Println(product)
	d := dbConn.Create(&product)
	if d.Error != nil {
		return nil, d.Error
	}

	return &product, nil
}

//GetProductByID **
func GetProductByID(ctx context.Context, id string) (*Product, error) {

	dbConn := NewClientConnection()
	defer dbConn.Close()

	product := Product{}
	err := dbConn.Where("id=?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

//UpdateProductStock **
func UpdateProductStock(ctx context.Context, updateStockReq UpdateStockReq) (*UpdateStockResp, error) {

	products := Product{}
	products.ID = uint(updateStockReq.ProductID)

	dbConn := NewClientConnection()
	defer dbConn.Close()

	if updateStockReq.RemoveStock != 0 {
		err := dbConn.Model(&products).Where("!qty < ?", updateStockReq.RemoveStock).Update("qty", gorm.Expr("qty - ?", updateStockReq.RemoveStock)).Error
		if err != nil {
			return nil, err
		}
	} else {

		err := dbConn.Model(&products).Where("id = ?", updateStockReq.ProductID).Update("qty", gorm.Expr("qty + ?", updateStockReq.AddStock)).Error
		if err != nil {
			return nil, err
		}

	}
	return &UpdateStockResp{
		Message: "Updated Successfully",
	}, nil
}

//GetProductList ***
func GetProductList(context.Context) ([]Product, error) {

	dbConn := NewClientConnection()
	defer dbConn.Close()

	products := []Product{}
	err := dbConn.Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil

}

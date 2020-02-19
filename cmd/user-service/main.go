package main

import (
	"log"
	"net/http"
	"time"

	"gitlab.com/scalent/ms-boilerplate/cmd/responses"
	"gitlab.com/scalent/ms-boilerplate/middleware"
	"gitlab.com/scalent/ms-boilerplate/models"
	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB
var JWTSecretKey = "secrect"
var err error

//LoginResponse Struct
type LoginResponse struct {
	Token          string    `json:"token,omitempty"`
	ExpirationTime time.Time `json:"-"`
	Error          string    `json:"error,omitempty"`
}

//Claims Struct
type Claims struct {
	Email    string `json:"email"`
	UserType string `json:"userType"`
	jwt.StandardClaims
}

func main() {

	DB, err = gorm.Open("mysql", "root:password@/ecommerce?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	router := mux.NewRouter()
	router.HandleFunc("/healthcheck", HealthCheck).Methods("GET")
	loginSubRouter := router.PathPrefix("/login").Subrouter()
	loginSubRouter.HandleFunc("/", UserLogin).Methods("POST", "OPTIONS")

	userSubRouter := router.PathPrefix("/user").Subrouter()
	// rm := middleware.RequestMiddleware{}
	userSubRouter.Use(middleware.ValidateMiddleware)
	userSubRouter.HandleFunc("/create", CreateUser).Methods("POST")
	userSubRouter.HandleFunc("/{id}", GetUserByID).Methods("GET")
	userSubRouter.HandleFunc("/", GetUsers).Methods("GET")
	defer DB.Close()
	log.Println("User Service:9001")
	log.Fatal(http.ListenAndServe(":9001", router))
}

//HealthCheck
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	responses.WriteOKResponse(w, "Health Check OK!")
}

//UserLogin Method
func UserLogin(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	email, password, _ := r.BasicAuth()
	input := models.UserLogin{
		Email:    email,
		Password: password,
	}
	loginResponse := LoginResponse{}
	user.Email = input.Email
	user.Password = input.Password

	err := DB.Table("users").Where("email=?", user.Email).First(&user).Error
	if err != nil {
		responses.WriteErrorResponse(w, http.StatusUnauthorized, "Unable To Login!")
		return
	}

	errp := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if errp != nil {
		responses.WriteErrorResponse(w, http.StatusUnauthorized, "Unable To Login!")
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, error := token.SignedString([]byte(JWTSecretKey))
	if error != nil {
		responses.WriteErrorResponse(w, http.StatusUnauthorized, "Unable To Login!")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	loginResponse.Token = tokenString
	responses.WriteOKResponse(w, loginResponse)
}

//CreateUser Method
func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	defer r.Body.Close()
	err := responses.ReadInput(r.Body, &user)
	if err != nil {
		responses.WriteErrorResponse(w, http.StatusBadRequest, "Invalid Input For Creating User!")
		return
	}
	pass, errp := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if errp != nil {
		responses.WriteErrorResponse(w, http.StatusBadRequest, "Invalid Password!")
		return
	}
	user.Password = string(pass)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err1 := DB.Table("users").Create(&user).Error
	if err1 != nil {
		responses.WriteErrorResponse(w, http.StatusBadRequest, "Unable To Create User!")
		return
	}
	responses.WriteOKResponse(w, "User Created Successfully!")
}

//GetUserByID Method
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]
	user := models.User{}
	err := DB.Table("users").Where("id=?", userID).First(&user).Error
	if err != nil {
		responses.WriteErrorResponse(w, http.StatusNotFound, "User Not Found!")
		return
	}
	responses.WriteOKResponse(w, user)
}

//GetUsers Method
func GetUsers(w http.ResponseWriter, r *http.Request) {
	users := []models.User{}
	err := DB.Table("users").Find(&users).Error
	if err != nil {
		responses.WriteErrorResponse(w, http.StatusNotFound, "No Users Found!")
		return
	}
	responses.WriteOKResponse(w, users)
}

package main

import (
	"net/http"

	"github.com/erply/api-go-wrapper/pkg/api"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Go Erply Middleware API
// @version 1.0
// @description Erply middleware swagger implementation in Go HTTP
// @contact.name Periyasamy Chinnasamy
// @contact.email persrc@gmail.com
// @securityDefinitions.apikey AuthKey
// @in header
// @name sessionkey
// @host localhost:8080
// @BasePath
func main() {

	erplyClients = make(map[string]*api.Client)

	http.HandleFunc("/customer/add", AddCustomer)
	http.HandleFunc("/customer/get", GetCustomer)
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	http.HandleFunc("/login", HandleAuth)
	http.ListenAndServe(":8080", nil)

}

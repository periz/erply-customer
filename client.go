package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"

	_ "customer/docs"

	"github.com/erply/api-go-wrapper/pkg/api"
)

var erplyClients map[string]*api.Client

type authParam struct {
	Username     string `json:username`
	Password     string `json:Password`
	Customercode string `json:customercode`
}

type addCustomerResp struct {
	Status string `json:status`
}

type authResp struct {
	Sessionkey string `sessionkey`
}

// customer godoc
// @Summary Add Customer details
// @Description Add Customer
// @Tags customer
// @ID customer-data
// @Accept  json
// @Produce  json
// @Param customer body []map[string]interface{} true "get customer filter"
// @Success 200 {object} main.addCustomerResp
// @Router /customer/add [post]
// @Security AuthKey
func AddCustomer(rw http.ResponseWriter, r *http.Request) {

	apiClient, err := GetClient(r.Header.Get("sessionkey"))

	if err {
		cli := apiClient.CustomerManager
		var customerArray []map[string]interface{}
		json.NewDecoder(r.Body).Decode(&customerArray)

		ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
		resp, _ := cli.SaveCustomerBulk(ctx, customerArray, map[string]string{})

		response := addCustomerResp{Status: resp.Status.ResponseStatus}
		data, _ := json.Marshal(response)
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(data)
	} else {
		resp := map[string]string{"errorMessage": "No client object found. Authentication required."}
		data, _ := json.Marshal(resp)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write(data)

	}
}

// customer godoc
// @Summary Get Customer details
// @Description Get Customer
// @Tags customer
// @ID customer-datainfo
// @Accept  json
// @Produce  json
// @Param customer body map[string]interface{} true "Get Customer details"
// @Success 200 {object} []map[string]interface{}
// @Router /customer/get [post]
// @Security AuthKey
func GetCustomer(rw http.ResponseWriter, r *http.Request) {

	apiClient, err := GetClient(r.Header.Get("sessionkey"))

	if err {
		cli := apiClient.CustomerManager

		var filters map[string]string
		json.NewDecoder(r.Body).Decode(&filters)

		ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
		response, _ := cli.GetCustomers(ctx, filters)

		data, _ := json.Marshal(response)
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(data)
	} else {
		resp := map[string]string{"errorMessage": "No client object found. Authentication required."}
		data, _ := json.Marshal(resp)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write(data)

	}

}

// authLogin godoc
// @Summary Auth Login
// @Description Auth Login
// @Tags auth
// @ID auth-login
// @Accept  json
// @Produce  json
// @Param authLogin body main.authParam true "Auth Login Input"
// @Success 200 {object} main.authResp
// @Router /login [post]
func HandleAuth(rw http.ResponseWriter, r *http.Request) {
	var param authParam
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&param)
	if err != nil {
		fmt.Println("Error parsing post data")
		rw.WriteHeader(http.StatusNoContent)
		rw.Write(nil)
	}

	username := flag.String("u", param.Username, "")
	password := flag.String("p", param.Password, "")
	clientCode := flag.String("cc", param.Customercode, "")

	flag.Parse()

	connectionTimeout := 60 * time.Second
	transport := &http.Transport{
		DisableKeepAlives:     true,
		TLSClientConfig:       &tls.Config{},
		ResponseHeaderTimeout: connectionTimeout,
	}
	httpCl := &http.Client{Transport: transport}

	clBldr := api.ClientBuilder{
		UserName:                 *username,
		Password:                 *password,
		ClientCode:               *clientCode,
		DefaultSessionLenSeconds: 6000, //default length is smaller than the wait time so the session will expire before the next attempt which should repeat the auth
		HttpCli:                  httpCl,
	}

	apiClient := clBldr.Build()
	sessionKey, _ := apiClient.GetSession()
	erplyClients[sessionKey] = apiClient

	resp := authResp{Sessionkey: sessionKey}
	data, _ := json.Marshal(resp)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(data)

}

func GetClient(sk string) (*api.Client, bool) {

	if client, isPresent := erplyClients[sk]; isPresent {
		return client, true
	} else {
		fmt.Println("No client object found. Relogin required.")
		return nil, false

	}
}

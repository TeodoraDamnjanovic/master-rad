package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/rpc"
	"strconv"
)

type RequestPayload struct {
	Action  string         `json:"action"`
	Auth    AuthPayload    `json:"auth,omitempty"`
	Log     LogPayload     `json:"log,omitempty"`
	User    UserPayload    `json:"user,omitempty"`
	Mail    MailPayload    `json:"mail,omitempty"`
	Order   OrderPayload   `json:"order,omitempty"`
	Payment PaymentPayload `json:"payment,omitempty"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type UserPayload struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Password  string `json:"password"`
}

type OrderPayload struct {
	UserId    int    `json:"user_id"`
	InvoiceId string `json:"invoice_id"`
	Paid      bool   `json:"paid"`
	Amount    int    `json:"amount"`
}

type PaymentPayload struct {
	UserId int     `json:"user_id"`
	Amount float64 `json:"amount"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "signup":
		app.signup(w, requestPayload.User)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	case "log":
		app.logItemViaRPC(w, requestPayload.Log)
	case "create-order":
		app.createOrder(w, requestPayload.Order)
	case "create-payment":
		app.createPayment(w, requestPayload.Payment)
	case "get-payments":
		app.getPayments(w, requestPayload.Payment)
	case "balance":
		app.getBalance(w, requestPayload.Payment)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("login - invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		log.Printf("login  %+v\n", a)
		app.errorJSON(w, errors.New("login - error calling auth service"))
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) signup(w http.ResponseWriter, a UserPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t")
	log.Printf("signup  %+v\n", a)
	request, err := http.NewRequest("POST", "http://authentication-service/signup", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		log.Printf("Signup user error %+v\n", response)
		app.errorJSON(w, errors.New("signup - error calling auth service"))
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Signed up"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)

}

func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	mailServiceURL := "http://mailer-service/send"

	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling mail service"))
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Mail has been sent to " + msg.To

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) createOrder(w http.ResponseWriter, o OrderPayload) {
	jsonData, _ := json.MarshalIndent(o, "", "\t")

	request, err := http.NewRequest("POST", "http://order-service/create-order", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		log.Printf("Create order request failed %+v\n", o)
		app.errorJSON(w, errors.New("create order - error calling order service"))
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Order created"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) createPayment(w http.ResponseWriter, p PaymentPayload) {
	jsonData, _ := json.MarshalIndent(p, "", "\t")

	request, err := http.NewRequest("POST", "http://payment-service/payment", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		log.Printf("Create payment request failed %+v\n", p)
		app.errorJSON(w, errors.New("create payment - error calling payment service"))
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Payment created"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) getBalance(w http.ResponseWriter, p PaymentPayload) {

	url := "http://payment-service/payment/" + strconv.Itoa(p.UserId) + "/balance"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		log.Printf("Get balance request failed %+v\n", p)
		app.errorJSON(w, errors.New("get balance - error calling apyment service"))
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Palance retreived"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) getPayments(w http.ResponseWriter, p PaymentPayload) {
	url := "http://payment-service/payment/" + strconv.Itoa(p.UserId)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		log.Printf("Get payments for user id %+v request failed \n", p.UserId)
		app.errorJSON(w, errors.New("get payments - error calling payment service"))
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Balance retreived"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)

}

type RPCPayload struct {
	Name string
	Data string
}

func (app *Config) logItemViaRPC(w http.ResponseWriter, l LogPayload) {
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	rpcPayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}

	var result string
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: result,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

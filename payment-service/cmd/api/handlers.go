package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"payment-service/data"
	"strconv"
	"strings"
)

func (app *Config) createPayment(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		UserId int     `json:"user_id"`
		Amount float64 `json:"amount"`
	}

	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		log.Printf("CreatePayment_readJSON %+v\n", err)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	log.Printf("request payload %+v\n", requestPayload)
	payment := data.Payment{
		UserId: requestPayload.UserId,
		Amount: requestPayload.Amount,
	}

	paymentId, err := app.Models.Payment.Insert(payment)

	if err != nil {
		log.Printf("CreatePayment_Insert %+v\n", err)
		app.errorJSON(w, errors.New("CreatePayment_Insert failed"), http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Payment with id %d is created!", paymentId),
		Data:    payment,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) getAllPaymentsForUserId(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.RequestURI)
	if err != nil {
		app.errorJSON(w, errors.New("GetAllPaymentsForUserId_InvalidRequestURI"), http.StatusBadRequest)
		return
	}

	userId, err := strconv.Atoi(getFirstParam(u.Path))
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	payments, err := app.Models.Payment.GetAllByUserId(userId)

	if err != nil {
		app.errorJSON(w, errors.New("UpdatePayment_GetPaymentByUserId failed"), http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Payments for user id %d are retreived!", userId),
		Data:    payments,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) viewBalanceForUserId(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.RequestURI)
	if err != nil {
		app.errorJSON(w, errors.New("ViewBalanceForUserId_InvalidRequestURI"), http.StatusBadRequest)
		return
	}
	fmt.Println(getFirstParam(u.Path))
	userId, err := strconv.Atoi(getFirstParam(u.Path))

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	balance, err := app.Models.Payment.GetBalanceByUserId(userId)

	if err != nil {
		log.Printf("ViewBalanceForUserId_GetBalanceByUserId %+v\n", err)
		app.errorJSON(w, errors.New("ViewBalanceForUserId_GetBalanceByUserId failed"), http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Balance for the user with id %d", userId),
		Data:    balance,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func getFirstParam(path string) (ps string) {
	res := strings.Split(path, "/")
	return res[2]
}

func getLastParam(path string) (ps string) {

	res := strings.Split(path, "/")
	return res[len(res)-1]
}

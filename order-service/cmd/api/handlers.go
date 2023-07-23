package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"order-service/data"
	"time"
)

type Order struct {
	UserId    int       `json:"user-id"`
	CreatedAt time.Time `json:"created-at"`
	InvoiceId string    `json:"invoice-id"`
	Paid      bool      `json:"paid"`
	Amount    int       `json:"amount"`
}

func (app *Config) createOrder(w http.ResponseWriter, r *http.Request) {
	var requestPayload Order

	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		log.Printf("CreateOrder_readJSON %+v\n", err)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	order := data.Order{
		UserId:    requestPayload.UserId,
		CreatedAt: requestPayload.CreatedAt,
		InvoiceId: requestPayload.InvoiceId,
		Paid:      requestPayload.Paid,
		Amount:    requestPayload.Amount,
	}

	orderId, err := app.Models.Order.Insert(order)

	if err != nil {
		log.Printf("CreateOrder_Insert %+v\n", err)
		app.errorJSON(w, errors.New("CreateOrder_Insert failed"), http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Order with id %s is created!", orderId),
		Data:    order,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) updateOrder(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		OrderId   int       `json:"order-id"`
		UserId    int       `json:"user-id"`
		CreatedAt time.Time `json:"created-at"`
		InvoiceId string    `json:"invoice-id"`
		Paid      bool      `json:"paid"`
		Amount    int       `json:"amount"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Printf("Signup_readJsonErr %+v\n", err)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	//get order by userId
	order, err := app.Models.Order.GetOneByUserId(requestPayload.UserId, requestPayload.OrderId)

	if err != nil {
		log.Printf("UpdateOrder_GetOrderByUserId %+v\n", err)
		app.errorJSON(w, errors.New("UpdateOrder_GetOrderByUserId failed"), http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Order with id %s is created!", order.ID),
		Data:    order,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

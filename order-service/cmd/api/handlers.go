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
	UserId    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	InvoiceId string    `json:"invoice_id"`
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

	// Create order
	order := data.Order{
		UserId:    requestPayload.UserId,
		CreatedAt: requestPayload.CreatedAt,
		InvoiceId: requestPayload.InvoiceId,
		Paid:      requestPayload.Paid,
		Amount:    requestPayload.Amount,
	}
	order.Status = "created"
	log.Printf("CreateOrder %+v\n", order)
	orderId, err := app.insertOrder(order)

	order.ID = orderId

	if err != nil {
		log.Printf("CreateOrder_Insert %+v\n", err)
		app.errorJSON(w, errors.New("CreateOrder_Insert failed"), http.StatusInternalServerError)
		return
	}

	// Process payment
	err = processPayment(order)

	if err != nil {
		log.Printf("CreateOrder_ProcessPayment failed %+v\n", err)
		// Compensation action: cancel order
		app.cancelOrder(order)
		http.Error(w, "Payment service failed", http.StatusInternalServerError)
		return
	}

	app.completeOrder(order)

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
		Message: fmt.Sprintf("Order with id %d is created!", order.ID),
		Data:    order,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) cancelOrder(order data.Order) error {
	// Update order table with status cancelled
	order.Status = "cancelled"

	log.Printf("CancelOrder %+v\n", order)

	err := app.Models.Order.UpdateOrderStatus(order.ID, order.Status, order.Paid)
	if err != nil {
		log.Printf("CancelOrder_Update %+v\n", err)
		return err
	}
	return nil
}

func (app *Config) completeOrder(order data.Order) error {
	// Update order table with status completed
	order.Status = "completed"
	order.Paid = true
	err := app.Models.Order.UpdateOrderStatus(order.ID, order.Status, order.Paid)
	if err != nil {
		log.Printf("CompleteOrder_Update %+v\n", err)
		return err
	}
	return nil
}

func (app *Config) insertOrder(order data.Order) (int, error) {
	orderId, err := app.Models.Order.Insert(order)

	return orderId, err
}

func (app *Config) updateOrderStatus(order data.Order) error {
	err := app.Models.Order.UpdateOrderStatus(order.ID, order.Status, order.Paid)
	if err != nil {
		log.Printf("UpdateOrderStatus_UpdateOrderStatus %+v\n", err)
		return err
	}
	return nil
}

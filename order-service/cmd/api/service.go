package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"order-service/data"
	"time"

	"github.com/cenkalti/backoff"
)

type Payment struct {
	UserId int `json:"user_id"`
	Amount int `json:"amount"`
}

func processPayment(order data.Order) error {
	payment := Payment{
		UserId: order.UserId,
		Amount: int(order.Amount),
	}
	body, _ := json.Marshal(payment)
	req, err := http.NewRequest("POST", "http://localhost:8083/payment", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	operation := func() error {
		resp, err := http.DefaultClient.Do(req)
		log.Printf("Payment response: %+v\n", resp)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			return errors.New(string(bodyBytes))
		}
		return nil
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = 30 * time.Second

	err = backoff.Retry(operation, expBackoff)
	if err != nil {
		log.Printf("Payment service failed after retries: %v", err)
		return err
	}
	return nil
}

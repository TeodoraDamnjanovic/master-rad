package main

import (
	"authentication-service/data"
	"authentication-service/event"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) Signup(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name,omitempty"`
		LastName  string `json:"last_name,omitempty"`
		Password  string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Printf("Signup_readJsonErr %+v\n", err)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user := data.User{
		Email:     requestPayload.Email,
		FirstName: requestPayload.FirstName,
		LastName:  requestPayload.LastName,
		Password:  requestPayload.Password,
		Active:    1,
	}

	userId, err := app.Models.User.Insert(user)
	if err != nil {
		app.errorJSON(w, errors.New("insert user failed"), http.StatusInternalServerError)
		return
	}

	app.sendMailEventViaRabbit(user)

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("%s is signed up", user.Email),
		Data:    userId,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

type MailMessage struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (app *Config) sendMailEventViaRabbit(user data.User) {
	msg := MailMessage{
		From:    "master@rad.com",
		To:      user.Email,
		Subject: "User registration service",
		Message: "Hi, " + user.FirstName + " you are successfuly signed up!",
	}
	err := app.pushToQueue(msg)
	if err != nil {
		log.Println("sendMailEventViaRabbitErr", err)
		return
	}
}

type Payload struct {
	Name string      `json:"name"`
	Data MailMessage `json:"data"`
}

// pushToQueue pushes a message into RabbitMQ
func (app *Config) pushToQueue(msg MailMessage) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		log.Println("Push to queue err ", err)
		return err
	}
	payload := Payload{
		Name: "send-mail",
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")
	err = emitter.Push(string(j), "mail.send")
	if err != nil {
		log.Println("Push to queue err ", err)
		return err
	}
	return nil
}

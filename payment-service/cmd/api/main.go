package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"payment-service/data"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting payment service")

	// connect to db
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}
	log.Println("Connected to Postgres!")

	//setup config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}
	//setup web server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	//start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Printf("OpenDb_Open %+v\n", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Printf("OpenDb_Ping %+v\n", err)
		return nil, err
	}
	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
			log.Printf("ConnectToDb_openDB %+v\n", err)
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}

// func connect() (*amqp.Connection, error) {
// 	var counts int64
// 	var backOff = 1 * time.Second
// 	var connection *amqp.Connection

// 	for {
// 		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
// 		if err != nil {
// 			fmt.Println("RabbitMQ not yet ready...")
// 			counts++
// 		} else {
// 			connection = c
// 			break
// 		}

// 		if counts > 5 {
// 			fmt.Println(err)
// 			return nil, err
// 		}

// 		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
// 		log.Println("backing off...")
// 		time.Sleep(backOff)
// 		continue
// 	}

// 	return connection, nil
// }

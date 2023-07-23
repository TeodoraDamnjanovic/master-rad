package data

import (
	"context"
	"database/sql"
	"log"
	"time"
)

const dbTimeout = time.Second * 3

var db *sql.DB

// New is the function used to create an instance of the data package. It returns the type
// Model, which embeds all the types we want to be available to our application.
func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		Order: Order{},
	}
}

// Models is the type for this package. Note that any model that is included as a member
// in this type is available to us throughout the application, anywhere that the
// app variable is used, provided that the model is also added in the New function.
type Models struct {
	Order Order
}

type Order struct {
	ID        int       `json:"id"`
	UserId    int       `json:"user-id"`
	CreatedAt time.Time `json:"created-at"`
	InvoiceId string    `json:"invoice-id"`
	Paid      bool      `json:"paid"`
	Amount    int       `json:"amount"`
}

// GetAll returns a slice of all orders for user id
func (o *Order) GetAllByUserId(userId int) ([]*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query :=
		` 	select id, user_id, created_at, invoice_id, paid, amount
			from orders 
			where o.user_id = $1`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*Order

	for rows.Next() {
		var order Order
		err := rows.Scan(
			&order.ID,
			&order.UserId,
			&order.CreatedAt,
			&order.InvoiceId,
			&order.Paid,
			&order.Amount,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		orders = append(orders, &order)
	}

	return orders, nil
}

// GetOne returns one order by id
func (o *Order) GetOneByUserId(userId int, orderId int) (*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `	select id, user_id, created_at, invoice_id, paid, amount
				from orders 
				where user_id = $1 and id = $2`

	var order Order
	row := db.QueryRowContext(ctx, query, userId, orderId)

	err := row.Scan(
		&order.ID,
		&order.UserId,
		&order.CreatedAt,
		&order.InvoiceId,
		&order.Paid,
		&order.Amount,
	)

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (o *Order) Insert(order Order) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID int
	stmt := `insert into orders (user_id, created_at, invoice_id, paid, amount)
		values ($1, $2, $3, $4, $5) returning id`

	err := db.QueryRowContext(ctx, stmt,
		&order.UserId,
		time.Now(),
		&order.InvoiceId,
		&order.Paid,
		&order.Amount,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

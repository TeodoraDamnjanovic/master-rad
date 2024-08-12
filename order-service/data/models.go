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
	UserId    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	InvoiceId string    `json:"invoice_id"`
	Paid      bool      `json:"paid"`
	Amount    int       `json:"amount"`
	Status    string    `json:"status"`
}

// GetAll returns a slice of all orders for user id
func (o *Order) GetAllByUserId(userId int) ([]*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query :=
		` 	select id, user_id, created_at, invoice_id, paid, amount
			from public.orders 
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
			&order.Status,
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

	query := `	select id, user_id, created_at, invoice_id, paid, amount, status
				from public.orders 
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
		&order.Status,
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
	stmt := `insert into public.orders (user_id, created_at, invoice_id, paid, amount, status)
		values ($1, $2, $3, $4, $5, $6) returning id`

	err := db.QueryRowContext(ctx, stmt,
		&order.UserId,
		time.Now(),
		&order.InvoiceId,
		&order.Paid,
		&order.Amount,
		&order.Status,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (o *Order) Update(order Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	log.Println("Updating order", order)

	stmt := `update public.orders set user_id = $1, created_at = $2, invoice_id = $3, paid = $4, amount = $5, status = $6
		where id = $7`

	_, err := db.ExecContext(ctx, stmt,
		&order.UserId,
		&order.CreatedAt,
		&order.InvoiceId,
		&order.Paid,
		&order.Amount,
		&order.Status,
		&order.ID,
	)

	if err != nil {
		return err
	}

	return nil

}

func (o *Order) UpdateOrderStatus(orderId int, status string, paid bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `UPDATE public.orders SET status = $1, paid = $2 WHERE id = $3`

	log.Printf("Executing query: %s with status: %s, paid: %t, orderId: %d", query, status, paid, orderId)

	result, err := db.ExecContext(ctx, query, status, paid, orderId)
	if err != nil {
		log.Println("Error updating order status:", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return err
	}

	if rowsAffected == 0 {
		log.Println("No rows were updated. Check if the orderId exists.")
		return sql.ErrNoRows
	}

	log.Printf("Order status updated successfully. Rows affected: %d", rowsAffected)

	return nil
}

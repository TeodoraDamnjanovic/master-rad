package data

import (
	"context"
	"database/sql"
	"log"
	"time"
)

const dbTimeout = time.Second * 3

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		Payment: Payment{},
	}
}

type Models struct {
	Payment Payment
}

type Payment struct {
	ID        int       `json:"id"`
	UserId    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Balance   float64   `json:"balance,omitempty"`
	Amount    float64   `json:"amount,omitempty"`
}

func (p *Payment) GetAllByUserId(userId int) ([]*Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query :=
		` 	select id, user_id, created_at, balance, amount
			from payments 
			where user_id = $1`

	rows, err := db.QueryContext(ctx, query, userId)
	if err != nil {
		log.Println("GetAllByUserId_QueryContext", err)
		return nil, err
	}
	defer rows.Close()

	var payments []*Payment

	for rows.Next() {
		var payment Payment
		err := rows.Scan(
			&payment.ID,
			&payment.UserId,
			&payment.CreatedAt,
			&payment.Balance,
			&payment.Amount,
		)
		if err != nil {
			log.Println("GetAllByUserId_ErrorScanning", err)
			return nil, err
		}

		payments = append(payments, &payment)
	}

	return payments, nil
}

func (p *Payment) GetOneByUserId(userId int, paymentId int) (*Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `	select id, user_id, created_at, balance, amount
				from payments 
				where user_id = $1 and id = $2`

	var payment Payment
	row := db.QueryRowContext(ctx, query, userId, paymentId)

	err := row.Scan(
		&payment.ID,
		&payment.UserId,
		&payment.CreatedAt,
		&payment.Balance,
		&payment.Amount,
	)

	if err != nil {
		log.Println("GetOneByUserId_ErrorScanning", err)
		return nil, err
	}

	return &payment, nil
}

func (p *Payment) Insert(payment Payment) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	balance, err := p.GetBalanceByUserId(payment.UserId)
	if err != nil {
		zero := 0.0
		balance = &zero
	}

	newBalance := *balance + payment.Amount

	var newID int
	stmt := `insert into payments (user_id, created_at, balance, amount)
		values ($1, $2, $3, $4) returning id`

	err = db.QueryRowContext(ctx, stmt,
		&payment.UserId,
		time.Now(),
		&newBalance,
		&payment.Amount,
	).Scan(&newID)

	if err != nil {
		log.Println("Insert_ErrorScanning", err)
		return 0, err
	}

	return newID, nil
}

func (p *Payment) GetBalanceByUserId(userId int) (*float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `	select balance
				from payments 
				where user_id = $1 order by created_at desc limit 1`

	var balance float64
	row := db.QueryRowContext(ctx, query, userId)

	err := row.Scan(
		&balance,
	)

	if err != nil {
		log.Println("GetBalanceByUserId_ErrorScanning", err, userId)
		return nil, err
	}

	return &balance, nil
}

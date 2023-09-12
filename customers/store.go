package customers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Customer struct {
	Id          string
	Email       string
	PhoneNumber string
	FirstName   string
	LastName    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CustomerStore struct {
	session *pgx.Conn
}

func NewCustomerStore(session *pgx.Conn) *CustomerStore {
	return &CustomerStore{session: session}
}

func (s *CustomerStore) GetById(ctx context.Context, id string) (*Customer, error) {
	customer := &Customer{}
	err := s.session.QueryRow(ctx, "select id, email, phone_number, first_name, last_name, created_at, updated_at from customers where id=$1", id).Scan(
		&customer.Id,
		&customer.Email,
		&customer.PhoneNumber,
		&customer.FirstName,
		&customer.LastName,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("cannot scan customer.\n\n%w", err)
	}

	return customer, nil
}

func (s *CustomerStore) GetByEmail(ctx context.Context, email string) (*Customer, error) {
	customer := &Customer{}
	err := s.session.QueryRow(ctx, "select id, email, phone_number, first_name, last_name, created_at, updated_at from customers where email=$1", email).Scan(
		&customer.Id,
		&customer.Email,
		&customer.PhoneNumber,
		&customer.FirstName,
		&customer.LastName,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot scan customer.\n\n%w", err)
	}

	return customer, nil
}

func (s *CustomerStore) GetByPhoneNumber(ctx context.Context, phone_number string) (*Customer, error) {
	customer := &Customer{}
	err := s.session.QueryRow(ctx, "select id, email, phone_number, first_name, last_name, created_at, updated_at from customers where phone_number=$1", phone_number).Scan(
		&customer.Id,
		&customer.Email,
		&customer.PhoneNumber,
		&customer.FirstName,
		&customer.LastName,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot scan customer.\n\n%w", err)
	}

	return customer, nil
}

func (s *CustomerStore) GetBySales(ctx context.Context, limit uint32) ([]Customer, error) {
	customers := []Customer{}
	rows, err := s.session.Query(ctx, "select id, email, phone_number, first_name, last_name, created_at, updated_at from customers order by (select sum(total_price) from transactions where customer_id=id) limit $1", limit)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("cannot scan customer.\n\n%w", err)
	}

	for rows.Next() {
		customer := Customer{}
		err := rows.Scan(
			&customer.Id,
			&customer.Email,
			&customer.PhoneNumber,
			&customer.FirstName,
			&customer.LastName,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}

	return customers, nil
}

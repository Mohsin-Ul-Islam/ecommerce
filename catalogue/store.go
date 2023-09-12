package catalogue

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Product struct {
	Id                string
	Title             string
	Description       string
	ImageUrl          string
	Price             float64
	AvailableQuantity uint32
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type CatalogueStore struct {
	session *pgx.Conn
}

func NewCatalogueStore(session *pgx.Conn) *CatalogueStore {
	return &CatalogueStore{session: session}
}

func (s *CatalogueStore) GetProductById(ctx context.Context, id string) (*Product, error) {
	product := &Product{}
	err := s.session.QueryRow(ctx, "select id, title, description, image_url, price, available_quantity, created_at, updated_at from products where id=$1", id).Scan(
		&product.Id,
		&product.Title,
		&product.Description,
		&product.ImageUrl,
		&product.Price,
		&product.AvailableQuantity,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("cannot scan product.\n\n%w", err)
	}

	return product, nil
}

func (s *CatalogueStore) AdjustProductQuantityById(ctx context.Context, id string, adjustment int32) (*Product, error) {
	product := &Product{}
	err := s.session.QueryRow(ctx, "update products set available_quantity = available_quantity + $1 where id = $2 returning id, title, description, image_url, price, available_quantity, created_at, updated_at", adjustment, id).Scan(
		&product.Id,
		&product.Title,
		&product.Description,
		&product.ImageUrl,
		&product.Price,
		&product.AvailableQuantity,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("cannot scan product.\n\n%w", err)
	}

	return product, nil
}

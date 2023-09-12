package transactions

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type TransactionStore struct {
	session *pgx.Conn
}

func NewTransactionStore(session *pgx.Conn) *TransactionStore {
	return &TransactionStore{session: session}
}

type TransactionModel struct {
	ID         string
	ProductID  string
	CustomerID string
	CreatedAt  time.Time
	Quantity   uint32
	TotalPrice float64
}

type ProductSales struct {
	ProductId  string
	TotalSales float64
}

func (s *TransactionStore) GetTransactionById(ctx context.Context, id string) (*TransactionModel, error) {
	var transactionModel = &TransactionModel{}

	err := s.session.QueryRow(ctx, "select id, product_id, customer_id, total_price, quantity, created_at from transactions where id = $1", id).Scan(
		&transactionModel.ID,
		&transactionModel.ProductID,
		&transactionModel.CustomerID,
		&transactionModel.TotalPrice,
		&transactionModel.Quantity,
		&transactionModel.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot scan transaction from database.\n\n%w", err)
	}

	return transactionModel, nil
}

func (s *TransactionStore) CreateTransaction(ctx context.Context, transaction TransactionModel) (*TransactionModel, error) {
	err := s.session.QueryRow(ctx, "insert into transactions (customer_id, product_id, quantity, total_price) values ($1, $2, $3, $4) returning id, created_at", transaction.CustomerID, transaction.ProductID, transaction.Quantity, transaction.TotalPrice).Scan(&transaction.ID, &transaction.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("cannot create transaction.\n\n%w", err)
	}

	return &transaction, nil
}

func (s *TransactionStore) GetTotalSales(ctx context.Context) (float64, error) {
	totalSales := float64(0)

	err := s.session.QueryRow(ctx, "select sum(total_price) as total_sales from transactions").Scan(&totalSales)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("cannot scan transaction from database.\n\n%w", err)
	}

	return totalSales, nil
}

func (s *TransactionStore) GetSalesByProduct(ctx context.Context) ([]ProductSales, error) {

	rows, err := s.session.Query(ctx, "select product_id, sum(total_price) as total_sales from transactions group by product_id")
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot scan transaction from database.\n\n%w", err)
	}

	productSalesList := []ProductSales{}
	for rows.Next() {
		productSales := ProductSales{}
		err := rows.Scan(&productSales.ProductId, &productSales.TotalSales)
		if err != nil {
			return nil, err
		}
		productSalesList = append(productSalesList, productSales)
	}

	return productSalesList, nil
}

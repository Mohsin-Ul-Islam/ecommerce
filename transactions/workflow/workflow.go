package workflow

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	cataloguepb_v1 "github.com/mohsin-ul-islam/ecommerce/catalogue/proto/v1"
	customerspb_v1 "github.com/mohsin-ul-islam/ecommerce/customers/proto/v1"
	"github.com/mohsin-ul-islam/ecommerce/transactions"
	pb "github.com/mohsin-ul-islam/ecommerce/transactions/proto/v1"
	"go.temporal.io/sdk/workflow"
)

type TransactionWorkflow struct {
	activities *WorkflowActivity
}

func (w *TransactionWorkflow) CreateTransactionWorkflow(ctx workflow.Context, request *pb.CreateTransactionRequest) (*pb.Transaction, error) {

	opts := workflow.ActivityOptions{StartToCloseTimeout: time.Second * 3}
	ctx = workflow.WithActivityOptions(ctx, opts)

	customer, err := w.activities.GetCustomerById(request.CustomerId)
	if err != nil {
		return nil, err
	}

	product, err := w.activities.GetProductById(request.ProductId)
	if err != nil {
		return nil, err
	}

	if product.AvailableQuantity != request.Quantity {
		return nil, errors.New("insufficient product quantity")
	}

	transaction := transactions.CreateTransactionRequestToModel(request)
	transaction.TotalPrice = float64(request.Quantity) * product.Price

	createdTransaction, err := w.activities.CreateTransaction(context.Background(), transaction)
	if err != nil {
		return nil, err
	}

	return transactions.TransactionToProto(createdTransaction), nil

}

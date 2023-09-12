package workflow

import (
	"context"
	cataloguepb_v1 "github.com/mohsin-ul-islam/ecommerce/catalogue/proto/v1"
	customerspb_v1 "github.com/mohsin-ul-islam/ecommerce/customers/proto/v1"
	"github.com/mohsin-ul-islam/ecommerce/transactions"
)

type WorkflowActivity struct {
	transactionStore *transactions.TransactionStore
	customersService customerspb_v1.CustomerServiceClient
	catalogueService cataloguepb_v1.CatalogueServiceClient
}

func (a *WorkflowActivity) GetCustomerById(ctx context.Context, id string) (*customerspb_v1.Customer, error) {
	return a.customersService.GetCustomerById(ctx, &customerspb_v1.GetCustomerByIdRequest{Id: id})
}

func (a *WorkflowActivity) GetProductById(ctx context.Context, id string) (*cataloguepb_v1.Product, error) {
	return &a.catalogueService.GetProductById(ctx, &cataloguepb_v1.GetProductByIdRequest{Id: id})
}

func (a *WorkflowActivity) CreateTransaction(ctx context.Context, transaction transactions.TransactionModel) (*transactions.TransactionModel, error) {
	return a.transactionStore.CreateTransaction(ctx, transaction)
}

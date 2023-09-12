package transactions

import (
	pb "github.com/mohsin-ul-islam/ecommerce/transactions/proto/v1"
	"time"
)

func TransactionToProto(t *TransactionModel) *pb.Transaction {
	return &pb.Transaction{
		Id:         t.ID,
		CustomerId: t.CustomerID,
		ProductId:  t.ProductID,
		Quantity:   t.Quantity,
		TotalPrice: t.TotalPrice,
		CreatedAt:  t.CreatedAt.String(),
	}
}

func CreateTransactionRequestToModel(r *pb.CreateTransactionRequest) TransactionModel {
	return TransactionModel{
		ID:         "",
		ProductID:  r.ProductId,
		CustomerID: r.CustomerId,
		CreatedAt:  time.Time{},
		Quantity:   r.Quantity,
		TotalPrice: 0,
	}
}

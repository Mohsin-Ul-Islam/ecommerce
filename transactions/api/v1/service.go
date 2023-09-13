package v1

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	cataloguepb_v1 "github.com/mohsin-ul-islam/ecommerce/catalogue/proto/v1"
	customerspb_v1 "github.com/mohsin-ul-islam/ecommerce/customers/proto/v1"
	"github.com/mohsin-ul-islam/ecommerce/transactions"
	pb "github.com/mohsin-ul-islam/ecommerce/transactions/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type TransactionService struct {
	pb.UnimplementedTransactionServiceServer
	transactionsStreamChan chan transactions.TransactionModel
	transactionStore       *transactions.TransactionStore
	customerService        customerspb_v1.CustomerServiceClient
	catalogueService       cataloguepb_v1.CatalogueServiceClient
}

func NewTransactionService(conn *pgx.Conn) *TransactionService {
	customerServiceConn, err := grpc.Dial("localhost:8083", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	catalogueServiceConn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	transactionStreamChan := make(chan transactions.TransactionModel)
	customerServiceClient := customerspb_v1.NewCustomerServiceClient(customerServiceConn)
	catalogueServiceClient := cataloguepb_v1.NewCatalogueServiceClient(catalogueServiceConn)
	return &TransactionService{transactionStore: transactions.NewTransactionStore(conn), customerService: customerServiceClient, catalogueService: catalogueServiceClient, transactionsStreamChan: transactionStreamChan}
}

func (s *TransactionService) GetTransactionById(ctx context.Context, request *pb.GetTransactionRequest) (
	*pb.Transaction,
	error,
) {
	transactionModel, err := s.transactionStore.GetTransactionById(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	if transactionModel == nil {
		return nil, status.Error(codes.NotFound, "transaction does not exist")
	}

	return transactions.TransactionToProto(transactionModel), nil
}

func (s *TransactionService) CreateTransaction(
	ctx context.Context,
	request *pb.CreateTransactionRequest,
) (*pb.Transaction, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	_, err := s.customerService.GetCustomerById(ctx, &customerspb_v1.GetCustomerByIdRequest{Id: request.CustomerId})
	if err != nil {
		return nil, err
	}

	product, err := s.catalogueService.GetProductById(ctx, &cataloguepb_v1.GetProductByIdRequest{Id: request.ProductId})
	if err != nil {
		return nil, err
	}

	if product.AvailableQuantity < request.Quantity {
		return nil, status.Error(codes.FailedPrecondition, "insufficient product available quantity")
	}

	transactionModel := transactions.CreateTransactionRequestToModel(request)
	transactionModel.TotalPrice = product.Price * float64(request.Quantity)
	createdTransaction, err := s.transactionStore.CreateTransaction(ctx, transactionModel)
	if err != nil {
		return nil, err
	}

	_, err = s.catalogueService.AdjustQuantity(ctx, &cataloguepb_v1.AdjustQuantityRequest{Id: product.Id, Adjustment: -int32(request.Quantity)})
	if err != nil {
		return nil, err
	}
	select {
	case s.transactionsStreamChan <- *createdTransaction:
	default:
	}
	return transactions.TransactionToProto(createdTransaction), nil
}

func (s *TransactionService) ListTransactions(
	request *pb.ListTransactionsRequest,
	stream pb.TransactionService_ListTransactionsServer,
) error {

	for {
		transaction := <-s.transactionsStreamChan
		err := stream.Send(transactions.TransactionToProto(&transaction))
		if err != nil {
			return fmt.Errorf("error in stream.Send.\n\n%w", err)
		}
	}
}

func (s *TransactionService) GetTotalSales(ctx context.Context, request *pb.GetTotalSalesRequest) (
	*pb.GetTotalSalesResponse,
	error,
) {
	totalSales, err := s.transactionStore.GetTotalSales(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.GetTotalSalesResponse{TotalSales: totalSales}, nil
}

func (s *TransactionService) GetSalesByProduct(ctx context.Context, request *pb.GetSalesByProductRequest) (
	*pb.GetSalesByProductResponse,
	error,
) {
	salesByProductList, err := s.transactionStore.GetSalesByProduct(ctx)
	if err != nil {
		return nil, err
	}

	salesByProductProtos := []*pb.SalesByProduct{}
	for _, salesByProduct := range salesByProductList {
		salesByProductProtos = append(salesByProductProtos, &pb.SalesByProduct{ProductId: salesByProduct.ProductId, TotalSales: salesByProduct.TotalSales})
	}

	return &pb.GetSalesByProductResponse{SalesByProducts: salesByProductProtos}, nil
}

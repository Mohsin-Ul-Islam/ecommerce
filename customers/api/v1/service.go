package v1

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/mohsin-ul-islam/ecommerce/customers"
	pb "github.com/mohsin-ul-islam/ecommerce/customers/proto/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CustomerService struct {
	pb.UnimplementedCustomerServiceServer
	customerStore *customers.CustomerStore
}

func NewCustomerService(conn *pgx.Conn) *CustomerService {
	return &CustomerService{customerStore: customers.NewCustomerStore(conn)}
}

func (s *CustomerService) GetCustomerById(ctx context.Context, request *pb.GetCustomerByIdRequest) (
	*pb.Customer,
	error,
) {
	customerModel, err := s.customerStore.GetById(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	if customerModel == nil {
		return nil, status.Error(codes.NotFound, "customer does not exist")
	}

	return customers.CustomerToProto(customerModel), nil
}

func (s *CustomerService) GetCustomersOrderedBySales(ctx context.Context, request *pb.GetCustomersOrderedBySalesRequest) (
	*pb.CustomersList,
	error,
) {
	customerModels, err := s.customerStore.GetBySales(ctx, request.Limit)
	if err != nil {
		return nil, err
	}

	customerProtos := []*pb.Customer{}
	for _, customer := range customerModels {
		customerProtos = append(customerProtos, customers.CustomerToProto(&customer))
	}

	return &pb.CustomersList{Customers: customerProtos}, nil
}

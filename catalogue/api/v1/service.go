package v1

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/mohsin-ul-islam/ecommerce/catalogue"
	pb "github.com/mohsin-ul-islam/ecommerce/catalogue/proto/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CatalogueService struct {
	pb.UnimplementedCatalogueServiceServer
	catalogueStore *catalogue.CatalogueStore
}

func NewCatalogueService(conn *pgx.Conn) *CatalogueService {
	return &CatalogueService{catalogueStore: catalogue.NewCatalogueStore(conn)}
}

func (s *CatalogueService) GetProductById(ctx context.Context, request *pb.GetProductByIdRequest) (
	*pb.Product,
	error,
) {
	productModel, err := s.catalogueStore.GetProductById(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	if productModel == nil {
		return nil, status.Error(codes.NotFound, "product does not exist")
	}

	return catalogue.ProductToProto(productModel), nil
}

func (s *CatalogueService) AdjustQuantity(ctx context.Context, request *pb.AdjustQuantityRequest) (
	*pb.Product,
	error) {

	productModel, err := s.catalogueStore.AdjustProductQuantityById(ctx, request.Id, request.Adjustment)
	if err != nil {
		return nil, err
	}

	if productModel == nil {
		return nil, status.Error(codes.NotFound, "product does not exist")
	}

	return catalogue.ProductToProto(productModel), nil

}

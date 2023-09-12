package catalogue

import pb "github.com/mohsin-ul-islam/ecommerce/catalogue/proto/v1"

func ProductToProto(c *Product) *pb.Product {
	return &pb.Product{
		Id:                c.Id,
		Title:             c.Title,
		Price:             c.Price,
		ImageUrl:          c.ImageUrl,
		Description:       c.Description,
		AvailableQuantity: c.AvailableQuantity,
	}
}

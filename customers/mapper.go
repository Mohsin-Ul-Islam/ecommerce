package customers

import pb "github.com/mohsin-ul-islam/ecommerce/customers/proto/v1"

func CustomerToProto(c *Customer) *pb.Customer {
	return &pb.Customer{
		Id:          c.Id,
		Email:       c.Email,
		PhoneNumber: c.PhoneNumber,
		FirstName:   c.FirstName,
		LastName:    c.LastName,
	}
}

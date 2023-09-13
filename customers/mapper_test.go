package customers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCustomerToProto(t *testing.T) {
	customer := Customer{
		Id:          "1",
		Email:       "example@mail.com",
		PhoneNumber: "+9661234567",
		FirstName:   "Foo",
		LastName:    "Bar",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}

	customerProto := CustomerToProto(&customer)
	assert.Equal(t, customer.Id, customerProto.Id)
	assert.Equal(t, customer.Email, customerProto.Email)
	assert.Equal(t, customer.PhoneNumber, customerProto.PhoneNumber)
	assert.Equal(t, customer.FirstName, customerProto.FirstName)
	assert.Equal(t, customer.LastName, customerProto.LastName)
}

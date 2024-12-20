package stripe

import (
	"github.com/stripe/stripe-go/v73"
	"github.com/stripe/stripe-go/v73/customer"
)

type CreateCustomerParams struct {
	Email       string
	FullName    string
	PhoneNumber string
}

func CreateCustomer(email string) (string, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
	}

	c, err := customer.New(params)
	if err != nil {
		return "", err
	}

	return c.ID, nil
}

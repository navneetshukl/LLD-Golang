package paypal

import "fmt"

type Paypal struct {
    email string
}

func NewPaypal(email string) *Paypal {
    return &Paypal{
        email: email,
    }
}

func (p *Paypal) Pay(amount float64) {
    fmt.Printf("Paid %.2f using Paypal: %s\n", amount, p.email)
}
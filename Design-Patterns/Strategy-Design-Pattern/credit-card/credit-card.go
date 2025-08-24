package creditcard

import "fmt"

type CreditCard struct {
    cardNumber string
    name       string
}

func NewCreditCard(name, cardNumber string) *CreditCard {
    return &CreditCard{
        cardNumber: cardNumber,
        name:       name,
    }
}

func (c *CreditCard) Pay(amount float64) {
    fmt.Printf("Paid %.2f using Credit Card (%s): %s\n", amount, c.name, c.cardNumber)
}
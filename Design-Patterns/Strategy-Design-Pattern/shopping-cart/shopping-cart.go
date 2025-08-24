package shoppingcart

import (
	"fmt"
	paymentstrategy "strategy-design/payment-strategy"
)

type ShoppingCart struct {
    payment paymentstrategy.PaymentStrategy
}

func NewShoppingCart(strategy paymentstrategy.PaymentStrategy) *ShoppingCart {
    return &ShoppingCart{
        payment: strategy,
    }
}

func (s *ShoppingCart) Checkout(amount float64) {
    if s.payment == nil {
        fmt.Println("Select the payment method first")
        return
    }
    s.payment.Pay(amount)
}

func (s *ShoppingCart) SetPaymentMethod(newMethod paymentstrategy.PaymentStrategy) {
    s.payment = newMethod
}
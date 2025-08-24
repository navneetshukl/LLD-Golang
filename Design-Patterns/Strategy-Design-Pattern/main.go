package main

import (
	"strategy-design/bitcoin"
	creditcard "strategy-design/credit-card"
	"strategy-design/paypal"
	shoppingcart "strategy-design/shopping-cart"
)

func main() {
	creditCardPayment := creditcard.NewCreditCard("MasterCard", "1234-5678-9012-3456")
	paypalPayment := paypal.NewPaypal("navneet@shukla.com")
	bitcoinPayment := bitcoin.NewBitcoin("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa")

	// Create shopping cart with credit card payment
	cart := shoppingcart.NewShoppingCart(creditCardPayment)
	cart.Checkout(123.45)

	// Switch to PayPal
	cart.SetPaymentMethod(paypalPayment)
	cart.Checkout(67.89)

	// Switch to Bitcoin
	cart.SetPaymentMethod(bitcoinPayment)
	cart.Checkout(999.99)

	// Demonstrate nil payment handling
	cart.SetPaymentMethod(nil)
	cart.Checkout(100.00)
}

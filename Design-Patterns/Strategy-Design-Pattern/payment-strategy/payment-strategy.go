package paymentstrategy

type PaymentStrategy interface{
	Pay(amount float64)
}
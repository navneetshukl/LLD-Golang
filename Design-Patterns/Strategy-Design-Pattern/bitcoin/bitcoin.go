package bitcoin

import "fmt"

type Bitcoin struct {
    walletAddress string
}

func NewBitcoin(wallet string) *Bitcoin {
    return &Bitcoin{
        walletAddress: wallet,
    }
}

func (b *Bitcoin) Pay(amount float64) {
    fmt.Printf("Paid %.2f using Bitcoin: %s\n", amount, b.walletAddress)
}
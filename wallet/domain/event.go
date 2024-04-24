package domain

type Event interface {
	isEvent()
}

func (w WalletCreated) isEvent()    {}
func (w WalletSpurious) isEvent()   {}
func (w WalletDeposited) isEvent()  {}
func (w WalletWithdrawed) isEvent() {}
func (w WalletWon) isEvent()        {}
func (w WalletLost) isEvent()       {}
func (w WalletReserved) isEvent()   {}
func (w WalletReleased) isEvent()   {}

type WalletCreated struct {
	ID int
}

type WalletSpurious struct {
	ID int
}

type WalletDeposited struct {
	ID     int
	Amount float64
}

type WalletWithdrawed struct {
	ID     int
	Amount float64
}

type WalletWon struct {
	ID     int
	Amount float64
}

type WalletLost struct {
	ID     int
	Amount float64
}

type WalletReserved struct {
	ID     int
	Amount float64
}

type WalletReleased struct {
	ID     int
	Amount float64
}

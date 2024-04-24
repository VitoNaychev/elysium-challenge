package domain

import (
	"errors"
)

type State byte

const (
	StateNew State = iota
	StateCreated
	StateSpurious
)

var (
	ErrInsufficientFunds     = errors.New("insufficient funds to perform operation")
	ErrUnsupportedTransition = errors.New("unsupported state transition")
	ErrStateSpurious         = errors.New("can't accept events while in spurious state")
)

type Wallet struct {
	id      int
	balance float64
	state   State

	changes []Event
	version int
}

func NewWallet() Wallet {
	return Wallet{
		id:      0,
		balance: 0,
		state:   StateNew,
	}
}

func NewWalletFromEvents(events []Event) Wallet {
	wallet := NewWallet()

	for _, event := range events {
		wallet.On(event, false)
	}

	return wallet
}

func (w *Wallet) GetID() int {
	return w.id
}

func (w *Wallet) GetBalance() float64 {
	return w.balance
}

func (w *Wallet) GetState() State {
	return w.state
}

func (w *Wallet) Create(userID int) error {
	if w.state != StateNew {
		return ErrUnsupportedTransition
	}

	w.raise(&WalletCreated{
		ID: userID,
	})
	return nil
}

func (w *Wallet) Deposit(amount float64) error {
	if w.state == StateSpurious {
		return ErrStateSpurious
	}

	w.raise(&WalletDeposited{
		ID:     w.id,
		Amount: amount,
	})
	return nil
}

func (w *Wallet) Withdraw(amount float64) error {
	if w.state == StateSpurious {
		return ErrStateSpurious
	}
	if w.balance-amount < 0 {
		return ErrInsufficientFunds
	}

	w.raise(&WalletWithdrawed{
		ID:     w.id,
		Amount: amount,
	})
	return nil
}

func (w *Wallet) Win(amount float64) error {
	if w.state == StateSpurious {
		return ErrStateSpurious
	}

	w.raise(&WalletWon{
		ID:     w.id,
		Amount: amount,
	})
	return nil
}

func (w *Wallet) Lose(amount float64) error {
	if w.state == StateSpurious {
		return ErrStateSpurious
	}

	if w.balance-amount < 0 {
		w.raise(&WalletSpurious{
			ID: w.id,
		})
		return ErrInsufficientFunds
	}

	w.raise(&WalletLost{
		ID:     w.id,
		Amount: amount,
	})
	return nil
}

func (w *Wallet) Reserve(amount float64) error {
	if w.state == StateSpurious {
		return ErrStateSpurious
	}
	if w.balance-amount < 0 {
		return ErrInsufficientFunds
	}

	w.raise(&WalletReserved{
		ID:     w.id,
		Amount: amount,
	})
	return nil
}

func (w *Wallet) Release(amount float64) error {
	if w.state == StateSpurious {
		return ErrStateSpurious
	}

	w.raise(&WalletReleased{
		ID:     w.id,
		Amount: amount,
	})
	return nil
}

func (w *Wallet) On(event Event, new bool) {
	switch e := event.(type) {
	case *WalletCreated:
		w.id = e.ID
		w.balance = 0
		w.state = StateCreated
	case *WalletSpurious:
		w.state = StateSpurious
	case *WalletDeposited:
		w.balance += e.Amount
	case *WalletWithdrawed:
		w.balance -= e.Amount
	case *WalletWon:
		w.balance += e.Amount
	case *WalletLost:
		w.balance -= e.Amount
	case *WalletReserved:
		w.balance -= e.Amount
	case *WalletReleased:
		w.balance += e.Amount
	}

	if !new {
		w.version++
	}
}

func (w *Wallet) Events() []Event {
	return w.changes
}

func (w *Wallet) Version() int {
	return w.version
}

func (w *Wallet) raise(event Event) {
	w.changes = append(w.changes, event)
	w.On(event, true)
}

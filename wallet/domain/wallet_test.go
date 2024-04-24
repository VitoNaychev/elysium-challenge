package domain_test

import (
	"testing"

	"github.com/VitoNaychev/elysium-challenge/assert"
	"github.com/VitoNaychev/elysium-challenge/wallet/domain"
)

func TestWalletConstructors(t *testing.T) {
	t.Run("creates new wallet and sets balance to zero and state to StateOK", func(t *testing.T) {
		wantBalance := float64(0)
		wantState := domain.StateNew

		wallet := domain.NewWallet()
		assert.Equal(t, wallet.GetBalance(), wantBalance)
		assert.Equal(t, wallet.GetState(), wantState)
	})

	t.Run("reconstructs wallet state from array of events", func(t *testing.T) {
		createdEvent := &domain.WalletDeposited{
			ID: 12,
		}
		depositedEvent := &domain.WalletDeposited{
			Amount: 100.00,
		}
		reservedEvent := &domain.WalletReserved{
			Amount: 50.00,
		}
		releasedEvent := &domain.WalletReleased{
			Amount: 50.00,
		}
		wonEvent := &domain.WalletWon{
			Amount: 50.00,
		}
		withdrawedEvent := &domain.WalletWithdrawed{
			Amount: 75.00,
		}

		wantBalance := depositedEvent.Amount -
			reservedEvent.Amount +
			releasedEvent.Amount +
			wonEvent.Amount -
			withdrawedEvent.Amount

		events := []domain.Event{
			createdEvent,
			depositedEvent,
			reservedEvent,
			releasedEvent,
			wonEvent,
			withdrawedEvent,
		}

		wallet := domain.NewWalletFromEvents(events)
		assert.Equal(t, wallet.GetBalance(), wantBalance)
	})
}

func TestWalletCreate(t *testing.T) {
	t.Run("saves create event and sets wallet ID", func(t *testing.T) {
		userID := 12

		wallet := domain.NewWallet()

		err := wallet.Create(userID)

		assert.RequireNoError(t, err)
		assert.Equal(t, wallet.GetID(), userID)
	})

	t.Run("returns ErrUnsupportedTransition on attempt to multiple Create calls", func(t *testing.T) {
		userID := 12

		wallet := domain.NewWallet()

		err := wallet.Create(userID)
		assert.RequireNoError(t, err)

		err = wallet.Create(userID)
		assert.Equal(t, err, domain.ErrUnsupportedTransition)
	})
}

func TestWalletDepositWithdraw(t *testing.T) {
	t.Run("saves deposit event and increases balance by set amount", func(t *testing.T) {
		userID := 12
		depositAmount := float64(100.99)

		wallet := createWallet(t, userID)

		wallet.Deposit(depositAmount)
		assert.Equal(t, wallet.GetBalance(), depositAmount)

		gotEventsCount := len(wallet.Events())
		wantEventsCount := 2

		requireEventsCount(t, gotEventsCount, wantEventsCount)
		assert.Type[*domain.WalletDeposited](t, wallet.Events()[gotEventsCount-1])
	})

	t.Run("returns ErrStateSpurious on deposit when wallet is in StateSpurious", func(t *testing.T) {
		userID := 12
		depositAmount := float64(100.99)
		loseAmount := float64(200.01)

		wallet := createWallet(t, userID)

		wallet.Lose(loseAmount)

		err := wallet.Deposit(depositAmount)
		assert.Equal(t, err, domain.ErrStateSpurious)
	})

	t.Run("saves withdraw event and decreases balance by set amount", func(t *testing.T) {
		userID := 12
		depositAmount := float64(100.99)
		withdrawAmount := float64(45.01)

		wallet := createWalletAndDeposit(t, userID, depositAmount)

		wallet.Withdraw(withdrawAmount)
		assert.Equal(t, wallet.GetBalance(), depositAmount-withdrawAmount)

		gotEventsCount := len(wallet.Events())
		wantEventsCount := 3

		requireEventsCount(t, gotEventsCount, wantEventsCount)
		assert.Type[*domain.WalletWithdrawed](t, wallet.Events()[gotEventsCount-1])
	})

	t.Run("returns ErrInsufficientFunds on insufficient funds during withdraw", func(t *testing.T) {
		userID := 12
		depositAmount := float64(100.99)
		withdrawAmount := float64(200.01)

		wallet := createWalletAndDeposit(t, userID, depositAmount)

		err := wallet.Withdraw(withdrawAmount)
		assert.Equal(t, err, domain.ErrInsufficientFunds)

		gotEventsCount := len(wallet.Events())
		wantEventsCount := 2

		requireEventsCount(t, gotEventsCount, wantEventsCount)
	})

	t.Run("returns ErrStateSpurious on withdraw when wallet is in StateSpurious", func(t *testing.T) {
		userID := 12
		withdrawAmount := float64(200.01)
		loseAmount := float64(200.01)

		wallet := createWallet(t, userID)

		wallet.Lose(loseAmount)

		err := wallet.Withdraw(withdrawAmount)
		assert.Equal(t, err, domain.ErrStateSpurious)
	})
}

func TestWalletWinLose(t *testing.T) {
	t.Run("saves win event and increases balance by set amount", func(t *testing.T) {
		userID := 12
		winAmount := float64(100.99)

		wallet := createWallet(t, userID)

		err := wallet.Win(winAmount)
		assert.RequireNoError(t, err)

		assert.Equal(t, wallet.GetBalance(), winAmount)

		gotEventsCount := len(wallet.Events())
		wantEventsCount := 2

		requireEventsCount(t, gotEventsCount, wantEventsCount)
		assert.Type[*domain.WalletWon](t, wallet.Events()[gotEventsCount-1])
	})

	t.Run("returns ErrStateSpurious on win when wallet is in StateSpurious", func(t *testing.T) {
		userID := 12
		winAmount := float64(100.99)
		loseAmount := float64(200.01)

		wallet := createWallet(t, userID)

		wallet.Lose(loseAmount)

		err := wallet.Win(winAmount)
		assert.Equal(t, err, domain.ErrStateSpurious)
	})

	t.Run("saves lost event and decreases balance by set amount", func(t *testing.T) {
		userID := 12
		depositAmount := float64(100.99)
		loseAmount := float64(45.01)

		wallet := createWalletAndDeposit(t, userID, depositAmount)

		wallet.Lose(loseAmount)
		assert.Equal(t, wallet.GetBalance(), depositAmount-loseAmount)

		gotEventsCount := len(wallet.Events())
		wantEventsCount := 3

		requireEventsCount(t, gotEventsCount, wantEventsCount)
		assert.Type[*domain.WalletLost](t, wallet.Events()[gotEventsCount-1])
	})

	t.Run("returns ErrInsufficientFunds on insufficient funds on lose", func(t *testing.T) {
		userID := 12
		depositAmount := float64(100.99)
		loseAmount := float64(200.01)

		wallet := createWalletAndDeposit(t, userID, depositAmount)

		err := wallet.Lose(loseAmount)
		assert.Equal(t, err, domain.ErrInsufficientFunds)

		gotEventsCount := len(wallet.Events())
		wantEventsCount := 3
		requireEventsCount(t, gotEventsCount, wantEventsCount)
	})

	t.Run("returns ErrStateSpurious on lose when wallet is in StateSpurious", func(t *testing.T) {
		userID := 12
		loseAmount := float64(200.01)

		wallet := createWallet(t, userID)

		wallet.Lose(loseAmount)

		err := wallet.Lose(loseAmount)
		assert.Equal(t, err, domain.ErrStateSpurious)
	})

	t.Run("sets state to StateSpurious on insufficient funds during lose", func(t *testing.T) {
		userID := 12
		depositAmount := float64(100.99)
		loseAmount := float64(200.01)

		wallet := createWalletAndDeposit(t, userID, depositAmount)

		err := wallet.Lose(loseAmount)
		assert.Equal(t, err, domain.ErrInsufficientFunds)
		assert.Equal(t, wallet.GetState(), domain.StateSpurious)

		gotEventsCount := len(wallet.Events())
		wantEventsCount := 3

		requireEventsCount(t, gotEventsCount, wantEventsCount)
		assert.Type[*domain.WalletSpurious](t, wallet.Events()[gotEventsCount-1])
	})
}

func TestWalletReserveRelease(t *testing.T) {
	t.Run("saves release event and increases balance by set amount", func(t *testing.T) {
		userID := 12
		releaseAmount := float64(100.99)

		wallet := createWallet(t, userID)

		wallet.Release(releaseAmount)
		assert.Equal(t, wallet.GetBalance(), releaseAmount)

		gotEventsCount := len(wallet.Events())
		wantEventsCount := 2

		requireEventsCount(t, gotEventsCount, wantEventsCount)
		assert.Type[*domain.WalletReleased](t, wallet.Events()[gotEventsCount-1])
	})

	t.Run("returns ErrStateSpurious on release when wallet is in StateSpurious", func(t *testing.T) {
		userID := 12
		releaseAmount := float64(100.99)
		loseAmount := float64(200.01)

		wallet := createWallet(t, userID)

		wallet.Lose(loseAmount)

		err := wallet.Release(releaseAmount)
		assert.Equal(t, err, domain.ErrStateSpurious)
	})

	t.Run("saves reserve event and decreases balance by set amount", func(t *testing.T) {
		userID := 12
		depositAmount := float64(100.99)
		reserveAmount := float64(45.01)

		wallet := createWalletAndDeposit(t, userID, depositAmount)

		wallet.Reserve(reserveAmount)
		assert.Equal(t, wallet.GetBalance(), depositAmount-reserveAmount)

		gotEventsCount := len(wallet.Events())
		wantEventsCount := 3

		requireEventsCount(t, gotEventsCount, wantEventsCount)
		assert.Type[*domain.WalletReserved](t, wallet.Events()[gotEventsCount-1])
	})

	t.Run("returns ErrInsufficientFunds on insufficient funds during reserve", func(t *testing.T) {
		userID := 12
		depositAmount := float64(100.99)
		reserveAmount := float64(200.01)

		wallet := createWalletAndDeposit(t, userID, depositAmount)

		err := wallet.Reserve(reserveAmount)
		assert.Equal(t, err, domain.ErrInsufficientFunds)

		gotEventsCount := len(wallet.Events())
		wantEventsCount := 2
		requireEventsCount(t, gotEventsCount, wantEventsCount)
	})

	t.Run("returns ErrStateSpurious on reserve when wallet is in StateSpurious", func(t *testing.T) {
		userID := 12
		reserveAmount := float64(200.01)
		loseAmount := float64(200.01)

		wallet := createWallet(t, userID)

		wallet.Lose(loseAmount)

		err := wallet.Reserve(reserveAmount)
		assert.Equal(t, err, domain.ErrStateSpurious)
	})
}

func requireEventsCount(t testing.TB, gotEventsCount, wantEventsCount int) {
	t.Helper()

	if gotEventsCount != wantEventsCount {
		t.Fatalf("got %v events want %v", gotEventsCount, wantEventsCount)
	}
}

func createWallet(t testing.TB, userID int) domain.Wallet {
	t.Helper()

	wallet := domain.NewWallet()

	err := wallet.Create(userID)
	assert.RequireNoError(t, err)

	return wallet
}

func createWalletAndDeposit(t testing.TB, userID int, amount float64) domain.Wallet {
	t.Helper()

	wallet := domain.NewWallet()

	err := wallet.Create(userID)
	assert.RequireNoError(t, err)

	err = wallet.Deposit(amount)
	assert.RequireNoError(t, err)

	return wallet
}

package auth_model

import (
	"net/http"
	"strings"

	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
)

func CreateWalletService(req models.CreateWalletRequest, db postgresql.Databases) (models.WalletBalance, int, error) {
	var (
		wallet = models.WalletBalance{
			AccountID: int(req.AccountID),
			Currency:  strings.ToUpper(req.Currency),
			Available: req.Available,
		}
	)

	_, err := wallet.GetWalletBalanceByAccountIDAndCurrency(db.Auth)
	if err == nil {
		return wallet, http.StatusOK, nil
	}

	err = wallet.CreateWalletBalance(db.Auth)
	if err != nil {
		return models.WalletBalance{}, http.StatusInternalServerError, err
	}

	return wallet, http.StatusCreated, nil

}

func UpdateWalletBalanceService(req models.UpdateWalletRequest, db postgresql.Databases) (models.WalletBalance, int, error) {
	var (
		wallet = models.WalletBalance{ID: req.ID}
	)

	code, err := wallet.GetWalletBalanceByID(db.Auth)
	if err != nil {
		return models.WalletBalance{}, code, err
	}

	wallet.Available = req.Available
	err = wallet.Update(db.Auth)
	if err != nil {
		return wallet, http.StatusInternalServerError, err
	}

	return wallet, http.StatusOK, nil

}

func GetWalletByAccountIDAndCurrencyService(db postgresql.Databases, accountID int, currency string) (models.WalletBalance, int, error) {
	var (
		wallet = models.WalletBalance{
			AccountID: accountID,
			Currency:  currency,
		}
	)

	code, err := wallet.GetWalletBalanceByAccountIDAndCurrency(db.Auth)
	if err != nil {
		return models.WalletBalance{}, code, err
	}

	return wallet, http.StatusOK, nil

}

func GetWalletsByAccountIDAndCurrenciesService(db postgresql.Databases, accountID int, currencies []string) (map[string]models.WalletBalance, int, error) {
	var (
		wallet            = models.WalletBalance{AccountID: accountID}
		walletBalancesMap = map[string]models.WalletBalance{}
	)

	walletBalances, err := wallet.GetWalletBalancesByAccountIDAndCurrencies(db.Auth, currencies)
	if err != nil {
		return walletBalancesMap, http.StatusInternalServerError, err
	}

	for _, w := range walletBalances {
		walletBalancesMap[w.Currency] = w
	}

	return walletBalancesMap, http.StatusOK, nil
}

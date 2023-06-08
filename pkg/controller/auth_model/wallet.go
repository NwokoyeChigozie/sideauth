package auth_model

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/auth-ms/services/auth_model"
	"github.com/vesicash/auth-ms/utility"
)

func (base *Controller) CreateWallet(c *gin.Context) {
	var (
		req models.CreateWalletRequest
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = postgresql.ValidateRequest(req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	walletBalance, code, err := auth_model.CreateWalletService(req, base.Db)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusCreated, "successful", walletBalance)
	c.JSON(http.StatusCreated, rd)

}

func (base *Controller) UpdateWalletBalance(c *gin.Context) {
	var (
		req models.UpdateWalletRequest
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = postgresql.ValidateRequest(req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	walletBalance, code, err := auth_model.UpdateWalletBalanceService(req, base.Db)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", walletBalance)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) GetWalletByAccountIDAndCurrency(c *gin.Context) {
	var (
		accountIDString = c.Param("account_id")
		currency        = strings.ToUpper(c.Param("currency"))
	)

	accountID, err := strconv.Atoi(accountIDString)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "invalid account id", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	walletBalance, code, err := auth_model.GetWalletByAccountIDAndCurrencyService(base.Db, accountID, currency)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", walletBalance)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) GetWalletsByAccountIDAndCurrencies(c *gin.Context) {
	var (
		accountIDString = c.Param("account_id")
		req             models.GetWalletsRequest
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = postgresql.ValidateRequest(req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	accountID, err := strconv.Atoi(accountIDString)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "invalid account id", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	walletBalances, code, err := auth_model.GetWalletsByAccountIDAndCurrenciesService(base.Db, accountID, req.Currencies)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", walletBalances)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) CreateWalletHistory(c *gin.Context) {
	var (
		req models.CreateWalletHistoryRequest
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = postgresql.ValidateRequest(req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	walletHistory := models.WalletHistory{
		AccountID:        strconv.Itoa(req.AccountID),
		Reference:        req.Reference,
		Amount:           req.Amount,
		Currency:         req.Currency,
		Type:             req.Type,
		AvailableBalance: req.AvailableBalance,
	}

	err = walletHistory.CreateWalletHistory(base.Db.Auth)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusCreated, "successful", walletHistory)
	c.JSON(http.StatusCreated, rd)

}

func (base *Controller) CreateWalletTransaction(c *gin.Context) {
	var (
		req models.CreateWalletTransactionRequest
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = postgresql.ValidateRequest(req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	walletTransaction := models.WalletTransaction{
		SenderAccountID:   strconv.Itoa(req.SenderAccountID),
		ReceiverAccountID: strconv.Itoa(req.ReceiverAccountID),
		SenderAmount:      req.SenderAmount,
		ReceiverAmount:    req.ReceiverAmount,
		SenderCurrency:    req.SenderCurrency,
		ReceiverCurrency:  req.ReceiverCurrency,
		Approved:          req.Approved,
		FirstApproval:     req.FirstApproval,
	}

	if req.SecondApproval != nil {
		walletTransaction.SecondApproval = *req.SecondApproval
	}

	err = walletTransaction.CreateWalletTransaction(base.Db.Auth)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusCreated, "successful", walletTransaction)
	c.JSON(http.StatusCreated, rd)

}

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

// CreateWallet
// GetWalletByAccountIDAndCurrency
// UpdateWalletBalance
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

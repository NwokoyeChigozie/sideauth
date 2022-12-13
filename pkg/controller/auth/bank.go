package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/auth-ms/services/auth"
	"github.com/vesicash/auth-ms/utility"
)

func (base *Controller) AddBankDetails(c *gin.Context) {
	var (
		req models.CreateBankRequest
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

	if models.MyIdentity.AccountID != req.AccountID {
		err := fmt.Errorf("not authorized to create bank detail for this user")
		rd := utility.BuildErrorResponse(http.StatusUnauthorized, "error", err.Error(), err, nil)
		c.JSON(http.StatusUnauthorized, rd)
		return
	}

	err = postgresql.ValidateRequest(req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	bankDetail, code, err := auth.CreateBankDetailService(req, base.Db)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusCreated, "bank details added", bankDetail)
	c.JSON(http.StatusCreated, rd)

}

func (base *Controller) GetBusinessCustomersBankDetails(c *gin.Context) {
	data, code, err := auth.GetBusinessCustomersBankDetailsService(base.Db, models.MyIdentity.AccountID)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "success", data)
	c.JSON(http.StatusOK, rd)
}

package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/auth-ms/services/auth"
	"github.com/vesicash/auth-ms/utility"
)

func (base *Controller) RequestPasswordReset(c *gin.Context) {
	var (
		req struct {
			EmailAddress string `json:"email_address"  pgvalidate:"exists=auth$users$email_address, email"`
			PhoneNumber  string `json:"phone_number"`
		}
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

	accountID, code, err := auth.RequestPasswordResetService(base.Db, req.EmailAddress, req.PhoneNumber)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Request Sent", gin.H{"account_id": accountID})
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) UpdatePasswordWithToken(c *gin.Context) {
	var (
		req struct {
			AccountID int    `json:"account_id" validate:"required" pgvalidate:"exists=auth$users$account_id"`
			Token     int    `json:"token" validate:"required"`
			Password  string `json:"password" validate:"required"`
		}
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

	code, err := auth.UpdatePasswordWithTokenService(base.Db, req.AccountID, req.Token, req.Password)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Password Updated", nil)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) UpdatePassword(c *gin.Context) {
	var (
		req struct {
			OldPassword string `json:"old_password" validate:"required"`
			NewPassword string `json:"new_password" validate:"required"`
		}
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

	code, err := auth.UpdatePassword(base.Db, models.MyIdentity.AccountID, req.OldPassword, req.NewPassword)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Password Updated", nil)
	c.JSON(http.StatusOK, rd)

}

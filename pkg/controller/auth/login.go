package auth

import (
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/services/auth"
	"github.com/vesicash/auth-ms/utility"
)

func (base *Controller) Login(c *gin.Context) {
	var (
		req models.LoginUserRequestModel
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

	data, code, err := auth.LoginService(c, base.Logger, req, base.Db)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "login successful", data)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) Logout(c *gin.Context) {
	user := models.User{AccountID: uint(models.MyIdentity.AccountID)}
	code, err := user.GetUserByAccountID(base.Db.Auth)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	user.LoginAccessToken = ""
	user.LoginAccessTokenExpiresIn = strconv.Itoa(int(time.Now().Unix()))
	err = user.Update(base.Db.Auth)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", err.Error(), err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "logout successful", nil)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) ValidateToken(c *gin.Context) {
	rd := utility.BuildSuccessResponse(http.StatusOK, "token valid", nil)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) PhoneOtpLogin(c *gin.Context) {
	var (
		req struct {
			PhoneNumber string `json:"phone_number" validate:"required"`
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

	accountID, code, err := auth.PhoneOtpLogin(c, base.Logger, req.PhoneNumber, base.Db)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "login successful", gin.H{"account_id": accountID})
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) GetAccessToken(c *gin.Context) {

	accessToken, code, err := auth.IssueAccessTokenService(base.Db, models.MyIdentity.AccountID)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "success", accessToken)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) ToggleMorStatus(c *gin.Context) {
	var (
		request models.EnableMORReq
	)

	err := c.ShouldBind(&request)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&request)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = postgresql.ValidateRequest(request)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	_, code, err := auth.UpdateUserMorSettings(base.Db, request, models.MyIdentity.AccountID)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "mor status updated!", nil)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) RevokeTokenHandler(c *gin.Context) {
	_, code, err := auth.RevokeAccessTokenService(base.Db, models.MyIdentity.AccountID)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Business access token has been revoked", nil)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) GetUserWalletBalance(c *gin.Context) {
	data, code, err := auth.GetUserWalletBalanceService(base.Db, models.MyIdentity.AccountID)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Data retrieved", data)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) GetUsers(c *gin.Context) {
	var (
		searchParam       = c.Query("search")
		isMorEnabledParam = c.Query("is_mor_enabled")
	)

	users, code, err := auth.GetUserService(base.Db, searchParam, isMorEnabledParam)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Users retrieved", users)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) ListSelectedCountries(c *gin.Context) {

	countries, code, err := auth.ListSelectedCountriesService(base.Db)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Countries retrieved", countries)
	c.JSON(http.StatusOK, rd)

}

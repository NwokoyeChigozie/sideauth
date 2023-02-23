package auth_model

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/services/auth_model"
	"github.com/vesicash/auth-ms/utility"
)

func (base *Controller) GetUser(c *gin.Context) {
	var (
		req models.GetUserModel
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

	user, code, err := auth_model.GetUserService(req, base.Db)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", gin.H{"user": user})
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) GetUsersByBusinessID(c *gin.Context) {
	var (
		businessID = c.Param("business_id")
	)
	fmt.Println(c.Params)

	businessIDint, err := strconv.Atoi(businessID)
	if err != nil {
		fmt.Println(err.Error())
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "incorrect business id type", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	users, code, err := auth_model.GetUsersByBusinessIDService(businessIDint, base.Db)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", users)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) SetAuthorizationRequired(c *gin.Context) {
	var (
		req struct {
			AccountID int  `json:"account_id" validate:"required" pgvalidate:"exists=auth$users$account_id"`
			Status    bool `json:"status" `
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

	user := models.User{AccountID: uint(req.AccountID)}
	code, err := user.GetUserByAccountID(base.Db.Auth)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}
	user.AuthorizationRequired = req.Status
	err = user.Update(base.Db.Auth)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", err.Error(), err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", true)
	c.JSON(http.StatusOK, rd)

}

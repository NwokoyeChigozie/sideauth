package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/services/auth"
	"github.com/vesicash/auth-ms/utility"
)

func (base *Controller) UpgradeUserTier(c *gin.Context) {
	var (
		req struct {
			Tier int `json:"tier" validate:"required,oneof=1 2"`
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

	code, err := auth.UpgradeUserTierService(base.Db, req.Tier, models.MyIdentity.AccountID)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Upgraded", nil)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) GetUserRestrictions(c *gin.Context) {

	data, code, err := auth.GetUserRestrictionsService(base.Logger, base.Db, models.MyIdentity.AccountID)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "success", data)
	c.JSON(http.StatusOK, rd)

}

package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/services/auth"
	"github.com/vesicash/auth-ms/utility"
)

func (base *Controller) UpgradeAccount(c *gin.Context) {
	var (
		req struct {
			BusinessType string `json:"business_type" validate:"required,oneof=ecommerce social_commerce marketplace"`
			BusinessName string `json:"business_name" validate:"required"`
			WebhookUri   string `json:"webhook_uri"`
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

	user, code, err := auth.UpgradeAccountService(base.Db, models.MyIdentity.AccountID, req.BusinessType, req.BusinessName, req.WebhookUri)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Upgraded", gin.H{"user": user})
	c.JSON(http.StatusOK, rd)

}

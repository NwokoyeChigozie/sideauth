package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/utility"
)

func (base *Controller) GetBusinessTypes(c *gin.Context) {
	var (
		businessType = models.BusinessType{}
	)

	businessTypes, err := businessType.GetBusinessTypes(base.Db.Auth)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", err.Error(), err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "data retrieved", businessTypes)
	c.JSON(http.StatusOK, rd)

}

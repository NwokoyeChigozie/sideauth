package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/services/auth"
	"github.com/vesicash/auth-ms/utility"
)

func (base *Controller) GetDisbursements(c *gin.Context) {

	data, code, err := auth.GetDisbursementsService(base.Db, models.MyIdentity.AccountID)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "success", data)
	c.JSON(http.StatusOK, rd)

}

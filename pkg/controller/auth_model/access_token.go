package auth_model

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/utility"
)

func (base *Controller) GetAccessToken(c *gin.Context) {
	var (
		accessToken = models.AccessToken{}
	)
	err := accessToken.GetAccessTokens(base.Db.Auth)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", err.Error(), err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", accessToken)
	c.JSON(http.StatusOK, rd)

}

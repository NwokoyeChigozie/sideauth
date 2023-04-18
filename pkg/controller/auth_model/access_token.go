package auth_model

import (
	"net/http"
	"strconv"

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
func (base *Controller) GetAccessTokenByKey(c *gin.Context) {
	var (
		key         = c.Param("key")
		accessToken = models.AccessToken{PrivateKey: key, PublicKey: key, IsLive: true}
	)

	code, err := accessToken.LiveTokensWithPublicOrPrivateKey(base.Db.Auth)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", accessToken)
	c.JSON(http.StatusOK, rd)
}
func (base *Controller) GetAccessTokenByBusinessID(c *gin.Context) {
	var (
		businessIDStr = c.Param("business_id")
		accessToken   = models.AccessToken{IsLive: true}
	)

	businessID, err := strconv.Atoi(businessIDStr)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "invalid business id type", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	accessToken.AccountID = businessID

	code, err := accessToken.GetLatestByAccountIDAndIsLive(base.Db.Auth)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", accessToken)
	c.JSON(http.StatusOK, rd)
}

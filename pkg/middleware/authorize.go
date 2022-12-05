package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/auth-ms/utility"
	"gorm.io/gorm"
)

const (
	ApiType       AuthorizationType = "api"
	AuthType      AuthorizationType = "auth"
	BusinessAdmin AuthorizationType = "business_admin"
	Business      AuthorizationType = "business"
)

type (
	AuthorizationType  string
	AuthorizationTypes []AuthorizationType
)

func Authorize(db postgresql.Databases, authTypes ...AuthorizationType) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(authTypes) > 0 {
			for _, v := range authTypes {
				msg, status := v.ValidateAuthorizationRequest(c, db)
				if !status {
					c.AbortWithStatusJSON(http.StatusUnauthorized, utility.UnauthorisedResponse(http.StatusUnauthorized, fmt.Sprint(http.StatusUnauthorized), "Unauthorized", msg))
					return
				}
			}
		}
	}
}

func (at AuthorizationType) in(authTypes AuthorizationTypes) bool {
	for _, v := range authTypes {
		if v == at {
			return true
		}
	}
	return false
}

func (at AuthorizationType) ValidateAuthorizationRequest(c *gin.Context, db postgresql.Databases) (string, bool) {
	if at == ApiType {
		return at.ValidateApiType(c, db)
	} else if at == AuthType {
		return "authorized", true // to be implemented
	} else if at == BusinessAdmin {
		return at.ValidateBusinessAdminType(c, db)
	} else if at == Business {
		return at.ValidateBusinessType(c, db)
	}

	return "authorized", true
}

func (at AuthorizationType) ValidateBusinessType(c *gin.Context, db postgresql.Databases) (string, bool) {
	_, msg, status := at.CheckAccessTokens(c, db)
	return msg, status
}
func (at AuthorizationType) ValidateBusinessAdminType(c *gin.Context, db postgresql.Databases) (string, bool) {
	token, msg, status := at.CheckAccessTokens(c, db)
	if !status {
		return msg, status
	}

	user := models.User{AccountID: uint(token.AccountID)}
	_, err := user.GetUserByAccountID(db.Auth)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "access denied", false
		} else {
			return "server error", false
		}
	}
	if user.AccountType != "admin" {
		return "access denied", false
	}
	return "authorized", true
}

func (at AuthorizationType) ValidateApiType(c *gin.Context, db postgresql.Databases) (string, bool) {
	_, msg, status := at.CheckAccessTokens(c, db)
	return msg, status
}

func (at AuthorizationType) CheckAccessTokens(c *gin.Context, db postgresql.Databases) (models.AccessToken, string, bool) {
	privateKey := GetHeader(c, "v-public-key")
	publicKey := GetHeader(c, "v-public-key")

	if privateKey == "" && publicKey == "" {
		return models.AccessToken{}, "missing api keys", false
	}

	if privateKey == "" || publicKey == "" {
		return models.AccessToken{}, "either public or private key is missing", false
	}

	token := models.AccessToken{PublicKey: publicKey, PrivateKey: privateKey, IsLive: true}
	_, err := token.LiveTokensWithPublicOrPrivateKey(db.Auth)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return token, "invalid keys", false
		} else {
			return token, "server error", false
		}
	}
	return token, "authorized", true
}

func GetHeader(c *gin.Context, key string) string {
	header := ""
	if c.GetHeader(key) != "" {
		header = c.GetHeader(key)
	} else if c.GetHeader(strings.ToLower(key)) != "" {
		header = c.GetHeader(strings.ToLower(key))
	} else if c.GetHeader(strings.ToUpper(key)) != "" {
		header = c.GetHeader(strings.ToUpper(key))
	} else if c.GetHeader(strings.Title(key)) != "" {
		header = c.GetHeader(strings.Title(key))
	}
	return header
}

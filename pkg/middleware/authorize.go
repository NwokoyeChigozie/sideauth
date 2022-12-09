package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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

			msg := ""
			for _, v := range authTypes {
				ms, status := v.ValidateAuthorizationRequest(c, db)
				if status {
					return
				}
				msg = ms
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, utility.UnauthorisedResponse(http.StatusUnauthorized, fmt.Sprint(http.StatusUnauthorized), "Unauthorized", msg))
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
		return at.ValidateAuthType(c, db)
	} else if at == BusinessAdmin {
		return at.ValidateBusinessAdminType(c, db)
	} else if at == Business {
		return at.ValidateBusinessType(c, db)
	}

	return "authorized", true
}

func (at AuthorizationType) ValidateAuthType(c *gin.Context, db postgresql.Databases) (string, bool) {

	var invalidToken = "Your request was made with invalid credentials."
	authorizationToken := GetHeader(c, "Authorization")
	if authorizationToken == "" {
		return "token not provided", false
	}

	bearerTokenArr := strings.Split(authorizationToken, " ")
	if len(bearerTokenArr) != 2 {
		return invalidToken, false
	}

	bearerToken := bearerTokenArr[1]

	if bearerToken == "" {
		return invalidToken, false
	}

	token, err := TokenValid(bearerToken)
	if err != nil {
		return invalidToken, false
	}

	claims := token.Claims.(jwt.MapClaims)
	activeUserType, ok := claims["type"].(string) //convert the interface to string
	if !ok {
		return invalidToken, false
	}

	activeUserAccountID, ok := claims["account_id"].(float64) //convert the interface to float
	if !ok {
		return invalidToken, false
	}

	authoriseStatus, ok := claims["authorised"].(bool) //check if token is authorised for middleware
	if !ok && !authoriseStatus {
		return invalidToken, false
	}

	myIdentity := models.UserIdentity{
		AccountID: int(activeUserAccountID),
		Type:      activeUserType,
	}

	user := models.User{AccountID: uint(myIdentity.AccountID)}
	code, err := user.GetUserByAccountID(db.Auth)
	if err != nil {
		if code == http.StatusInternalServerError {
			return err.Error(), false
		}
		return "user does not exist", false
	}

	if user.LoginAccessToken != bearerToken {
		return invalidToken, false
	}

	if user.LoginAccessTokenExpiresIn == "" {
		return invalidToken, false
	}

	parseInt, err := strconv.Atoi(user.LoginAccessTokenExpiresIn)
	if err != nil {
		return invalidToken, false
	}

	unixTimeUTC := time.Unix(int64(parseInt), 0)
	if time.Now().After(unixTimeUTC) {
		return "expired token", false
	}

	models.MyIdentity = &myIdentity
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

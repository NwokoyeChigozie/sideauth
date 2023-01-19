package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vesicash/auth-ms/pkg/controller/auth_model"
	"github.com/vesicash/auth-ms/pkg/middleware"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/auth-ms/utility"
)

func Model(r *gin.Engine, ApiVersion string, validator *validator.Validate, db postgresql.Databases, logger *utility.Logger) *gin.Engine {
	auth_model := auth_model.Controller{Db: db, Validator: validator, Logger: logger}

	modelTypeUrl := r.Group(fmt.Sprintf("%v/auth", ApiVersion), middleware.Authorize(db, middleware.AppType))
	{
		modelTypeUrl.POST("/get_user", auth_model.GetUser)
		modelTypeUrl.POST("/set_authorization_required", auth_model.SetAuthorizationRequired)
		modelTypeUrl.POST("/get_user_credentials", auth_model.GetUserCredentials)
		modelTypeUrl.POST("/create_user_credentials", auth_model.CreateUserCredentials)
		modelTypeUrl.POST("/update_user_credentials", auth_model.UpdateUserCredentials)
		modelTypeUrl.POST("/get_user_profile", auth_model.GetUserProfile)
		modelTypeUrl.POST("/get_country", auth_model.GetCountry)
		modelTypeUrl.POST("/get_bank_detail", auth_model.GetBankDetail)
		modelTypeUrl.POST("/get_business_profile", auth_model.GetBusinessProfile)
		modelTypeUrl.GET("/get_access_token", auth_model.GetAccessToken)
		modelTypeUrl.POST("/validate_on_db", auth_model.ValidateOnDB)
		modelTypeUrl.POST("/validate_authorization", auth_model.ValidateAuthorization)
		modelTypeUrl.POST("/get_authorize", auth_model.GetAuthorize)
		modelTypeUrl.POST("/create_authorize", auth_model.CreateAuthorize)
		modelTypeUrl.POST("/update_authorize", auth_model.UpdateAuthorize)

	}

	return r
}

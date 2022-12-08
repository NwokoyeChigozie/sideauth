package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vesicash/auth-ms/pkg/controller/auth"
	"github.com/vesicash/auth-ms/pkg/middleware"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
)

func Auth(r *gin.Engine, ApiVersion string, validator *validator.Validate, db postgresql.Databases) *gin.Engine {
	auth := auth.Controller{Db: db, Validator: validator}

	authUrl := r.Group(fmt.Sprintf("%v/auth", ApiVersion))
	{
		authUrl.POST("/signup", auth.Signup)
		authUrl.POST("/signup/bulk", auth.BulkSignup)
		authUrl.POST("/login", auth.Login)

	}

	authTypeUrl := r.Group(fmt.Sprintf("%v/auth", ApiVersion), middleware.Authorize(db, middleware.AuthType))
	{
		authTypeUrl.POST("/send_otp", auth.SendOTP)

	}
	authApiUrl := r.Group(fmt.Sprintf("%v/auth/api", ApiVersion), middleware.Authorize(db, middleware.ApiType))
	{
		authApiUrl.POST("/send_otp", auth.SendOTPAPI)

	}
	return r
}

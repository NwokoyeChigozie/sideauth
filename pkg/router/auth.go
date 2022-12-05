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

	}

	authUrl1 := r.Group(fmt.Sprintf("%v/auth", ApiVersion), middleware.Authorize(db, middleware.ApiType))
	{
		authUrl1.POST("/send_otp", auth.SendOTP)

	}
	return r
}

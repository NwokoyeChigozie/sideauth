package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vesicash/auth-ms/pkg/controller/auth"
	"github.com/vesicash/auth-ms/pkg/middleware"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/auth-ms/utility"
)

func Auth(r *gin.Engine, ApiVersion string, validator *validator.Validate, db postgresql.Databases, logger *utility.Logger) *gin.Engine {
	auth := auth.Controller{Db: db, Validator: validator, Logger: logger}

	authUrl := r.Group(fmt.Sprintf("%v", ApiVersion))
	{
		authUrl.POST("/signup", auth.Signup)
		authUrl.POST("/signup/bulk", auth.BulkSignup)

		authUrl.POST("/login", auth.Login)
		authUrl.POST("/login-phone", auth.PhoneOtpLogin)

		authUrl.POST("/otp/send_otp", auth.SendOTPAPI)
		authUrl.POST("/is_otp_valid", auth.ValidateOtp)

		authUrl.POST("/reset-password", auth.RequestPasswordReset)
		authUrl.POST("/reset-password/change-password", auth.UpdatePasswordWithToken)

		authUrl.POST("/contact-us", auth.ContactUs)
		authUrl.GET("/business-types", auth.GetBusinessTypes)

	}

	authTypeUrl := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db, middleware.AuthType))
	{
		authTypeUrl.POST("/send_otp", auth.SendOTP)

		authTypeUrl.POST("/user/bank_details", auth.AddBankDetails)
		authTypeUrl.POST("/user/update_tour_status", auth.UpdateTourStatus)

		authTypeUrl.GET("/user/restrictions", auth.GetUserRestrictions)
		authTypeUrl.POST("/user/upgrade_tier", auth.UpgradeUserTier)
		authTypeUrl.POST("/user/upgrade/account", auth.UpgradeAccount)

		authTypeUrl.POST("/user/security/update_password", auth.UpdatePassword)
		authTypeUrl.GET("/user/security/get_access_token", auth.GetAccessToken)

		authTypeUrl.GET("/user/disbursements", auth.GetDisbursements)

		authTypeUrl.GET("/business/customers/bank_details", auth.GetBusinessCustomersBankDetails)

		authTypeUrl.POST("/validate-token", auth.ValidateToken)
		authTypeUrl.POST("/logout", auth.Logout)

		authTypeUrl.POST("/revoke-token", auth.RevokeTokenHandler)

		authTypeUrl.GET("/account/wallet", auth.GetUserWalletBalance)

	}

	businessAdminUrl := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db, middleware.BusinessAdmin, middleware.AppType))
	{
		businessAdminUrl.POST("/toggle-mor-status", auth.ToggleMorStatus)
		businessAdminUrl.GET("/users/get", auth.GetUsers)

		businessAdminUrl.GET("/countries/mor", auth.ListSelectedCountries)

	}

	authApiUrl := r.Group(fmt.Sprintf("%v/api", ApiVersion), middleware.Authorize(db, middleware.ApiType))
	{
		authApiUrl.POST("/send_otp", auth.SendOTPAPI)

	}
	return r
}

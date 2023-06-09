package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vesicash/auth-ms/pkg/controller/health"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/auth-ms/utility"
)

func Health(r *gin.Engine, ApiVersion string, validator *validator.Validate, db postgresql.Databases, logger *utility.Logger) *gin.Engine {
	healthController := health.Controller{Db: db, Logger: logger}

	healthUrl := r.Group(fmt.Sprintf("%v", ApiVersion))
	{
		healthUrl.GET("/health", healthController.Get)
	}
	return r
}

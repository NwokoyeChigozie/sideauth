package health

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vesicash/auth-ms/internal/models"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/auth-ms/services/ping"
	"github.com/vesicash/auth-ms/utility"
)

type Controller struct {
	Db        postgresql.Databases
	Validator *validator.Validate
}

func (base *Controller) Post(c *gin.Context) {
	var (
		req    = models.Ping{}
		logger = utility.NewLogger()
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if !ping.ReturnTrue() {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "ping failed", fmt.Errorf("ping failed"), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	logger.Info("ping successfull")

	rd := utility.BuildSuccessResponse(http.StatusOK, "ping successful", req.Message)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) Get(c *gin.Context) {
	logger := utility.NewLogger()
	if !ping.ReturnTrue() {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "ping failed", fmt.Errorf("ping failed"), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	logger.Info("ping successfull")
	rd := utility.BuildSuccessResponse(http.StatusOK, "ping successful", gin.H{"user": "user object"})
	c.JSON(http.StatusOK, rd)

}

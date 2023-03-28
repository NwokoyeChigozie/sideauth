package status

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/auth-ms/utility"
)

type Controller struct {
	Db     postgresql.Databases
	Logger *utility.Logger
}

func (base *Controller) Get(c *gin.Context) {
	rd := utility.BuildSuccessResponse(http.StatusOK, "ping successful", gin.H{"user": "user object"})
	c.JSON(http.StatusOK, rd)

}

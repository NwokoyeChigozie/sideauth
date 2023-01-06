package auth

import (
	"github.com/go-playground/validator/v10"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/auth-ms/utility"
)

type Controller struct {
	Db        postgresql.Databases
	Validator *validator.Validate
	Logger    *utility.Logger
}

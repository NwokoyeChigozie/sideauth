package auth

import (
	"github.com/go-playground/validator/v10"
	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
)

type Controller struct {
	Db        postgresql.Databases
	Validator *validator.Validate
}

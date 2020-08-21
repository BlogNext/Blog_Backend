package validate

import (
	"github.com/go-playground/validator/v10"
	"github.com/blog_backend/exception"
)

type ValidateRequest interface {
	GetError(err validator.ValidationErrors) exception.MyException
}

package api

import (
	"github.com/DakshChawla/simplebank/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		util.IsSupportedCurrency(currency)
	}

	return true
}

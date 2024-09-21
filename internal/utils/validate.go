package utils

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// проверка не требуется, т.к. проверяется ключ и функция на nil, иначе - паника
		_ = v.RegisterValidation("luhn", luhnValidation)
	}
}

var luhnValidation validator.Func = func(fl validator.FieldLevel) bool {
	number, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	return LuhnCheck(number)
}

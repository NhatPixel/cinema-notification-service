package validator

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func RegisterNotificationValidation(v *validator.Validate) {
	v.RegisterValidation("future", func(fl validator.FieldLevel) bool {
		t, ok := fl.Field().Interface().(*time.Time)
		if !ok || t == nil {
			return true
		}
		return t.After(time.Now())
	})
}
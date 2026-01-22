package validation

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func TranslateValidationError(err error) error {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return err
	}

	for _, fe := range ve {
		switch fe.Field() {
		case "Title":
			switch fe.Tag() {
			case "required":
				return fmt.Errorf("title không được để trống")
			case "min":
				return fmt.Errorf("title phải có ít nhất %s ký tự", fe.Param())
			case "max":
				return fmt.Errorf("title không được vượt quá %s ký tự", fe.Param())
			}

		case "ExpiresAt":
			switch fe.Tag() {
			case "future":
				return fmt.Errorf("expires_at phải là thời gian trong tương lai")
			}
		}
	}

	return fmt.Errorf("dữ liệu không hợp lệ")
}

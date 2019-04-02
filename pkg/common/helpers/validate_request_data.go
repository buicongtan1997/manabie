package helpers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v8"
)

func ValidateBinding(c *gin.Context, out interface{}) (error) {
	err := c.ShouldBind(out)

	if err != nil {
		switch err.(type) {
		case validator.ValidationErrors:
			for _, e := range err.(validator.ValidationErrors) {
				switch e.Tag {
				case "required":
					return fmt.Errorf("%s is required", e.Name)
				case "max":
					return fmt.Errorf("%s must be less than %s", e.Field, e.Param)
				case "min":
					return fmt.Errorf("%s must be greater than %s", e.Field, e.Param)
				case "email":
					return fmt.Errorf("email invalid format")
				default:
					return fmt.Errorf("%s is invalid", e.Field)
				}
			}
		default:
			return err
		}
	}
	return nil
}


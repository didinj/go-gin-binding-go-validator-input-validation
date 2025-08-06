package controllers

import (
	"net/http"

	"github.com/didinj/go-input-validation/models"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var trans ut.Translator

func SetTranslator(t ut.Translator) {
	trans = t
}

func CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		var ve validator.ValidationErrors
		if ok := validator.As(err, &ve); ok {
			out := make(map[string]string)
			for _, fe := range ve {
				out[fe.Field()] = fe.Translate(trans)
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": out})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully!",
		"user":    user,
	})
}

func validationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "min":
		return fe.Field() + " must be at least " + fe.Param() + " characters"
	case "max":
		return fe.Field() + " must be at most " + fe.Param() + " characters"
	case "email":
		return "Invalid email address"
	case "gte":
		return fe.Field() + " must be greater than or equal to " + fe.Param()
	case "lte":
		return fe.Field() + " must be less than or equal to " + fe.Param()
	case "password":
		return "Password must be at least 8 characters and contain only letters, numbers, or symbols (!@#$%^&*)"
	default:
		return fe.Field() + " is invalid"
	}
}

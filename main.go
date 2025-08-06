package main

import (
	"github.com/didinj/go-input-validation/controllers"
	"github.com/didinj/go-input-validation/validators"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var trans ut.Translator

func main() {
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register custom password validation
		v.RegisterValidation("password", validators.PasswordValidator)

		v.RegisterTranslation("password", trans, func(ut ut.Translator) error {
			return ut.Add("password", "{0} must be at least 8 characters and contain only letters, numbers, or symbols", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("password", fe.Field())
			return t
		})

		// Set up English translator
		eng := en.New()
		uni := ut.New(eng, eng)
		trans, _ = uni.GetTranslator("en")

		// Register built-in translations
		en_translations.RegisterDefaultTranslations(v, trans)
	}

	controllers.SetTranslator(trans) // pass it to controller

	router.POST("/users", controllers.CreateUser)

	router.Run(":8080")
}

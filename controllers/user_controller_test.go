package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/didinj/go-input-validation/validators"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Validator setup
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("password", validators.PasswordValidator)

		eng := en.New()
		uni := ut.New(eng, eng)
		trans, _ := uni.GetTranslator("en")
		en_translations.RegisterDefaultTranslations(v, trans)

		v.RegisterTranslation("password", trans, func(ut ut.Translator) error {
			return ut.Add("password", "{0} must be at least 8 characters and contain only letters, numbers, or symbols", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("password", fe.Field())
			return t
		})

		SetTranslator(trans)
	}

	router.POST("/users", CreateUser)
	return router
}

func TestCreateUser_ValidInput(t *testing.T) {
	router := setupRouter()

	body := `{
		"name": "Alice",
		"email": "alice@example.com",
		"age": 25,
		"password": "Secure123!"
	}`

	req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "User created successfully")
}

func TestCreateUser_InvalidInput(t *testing.T) {
	router := setupRouter()

	body := `{
		"name": "Al",
		"email": "bad-email",
		"age": 10,
		"password": "123"
	}`

	req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Name must be at least 3 characters")
	assert.Contains(t, w.Body.String(), "Email must be a valid email address")
	assert.Contains(t, w.Body.String(), "Age must be greater than or equal to 18")
	assert.Contains(t, w.Body.String(), "Password must be at least 8 characters")
}

package config

import (
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/locales/ar"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"

	"github.com/joho/godotenv"
)

var Env = NewConfig()
var Validator = NewValidator()



type config struct {
	ServerPort string `json:"server_port"`

	DBHost string `json:"db_host"`
	DBPort string `json:"db_port"`
	DBName string `json:"db_name"`
	DBUser string `json:"db_user"`
	DBPass string `json:"db_pass"`
}

type validatorStruct struct {
	*validator.Validate
	*ut.UniversalTranslator
}

func NewConfig() *config {
	godotenv.Load()
	return &config{
		ServerPort: getEnv("SERVER_PORT", "8080"),

		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "27017"),
		DBName: getEnv("DB_NAME", "tubba"),
		DBUser: getEnv("DB_USER", ""),
		DBPass: getEnv("DB_PASS", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			return i
		}
	}
	return fallback
}


func NewValidator() *validatorStruct{
	en := en.New()
	ar := ar.New()
	uni := ut.New(en, en, ar)

	en_trans, _ := uni.GetTranslator("en")

	validator := validator.New()
	en_translations.RegisterDefaultTranslations(validator, en_trans)

	validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &validatorStruct{
		Validate: validator,
		UniversalTranslator: uni,
	}
}
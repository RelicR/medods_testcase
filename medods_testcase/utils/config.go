package utils

import (
	"os"
	"strconv"
)

type Config struct {
	DbHost       string
	DbPort       int
	DbUser       string
	DbPass       string
	DbName       string
	AppHost      string
	AppPort      int
	AppRefSecret string
	AppAcsSecret string
}

func NewConfig() *Config {
	return &Config{
		DbHost:       fromEnv("DB_HOST", "db"),
		DbPort:       fromEnvInt("DB_PORT", 5432),
		DbUser:       fromEnv("DB_USER", "testcase_api"),
		DbPass:       fromEnv("DB_PASSWORD", "testcase_api_password"),
		DbName:       fromEnv("DB_NAME", "auth_testcase"),
		AppHost:      fromEnv("APP_HOST", "0.0.0.0"),
		AppPort:      fromEnvInt("APP_PORT", 8080),
		AppRefSecret: fromEnv("APP_REFRESH_SECRET", "asdkljh23"),
		AppAcsSecret: fromEnv("APP_ACCESS_SECRET", "sdklsdf934"),
	}
}

func fromEnv(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func fromEnvInt(key string, defaultValue int) int {
	if value, ok := os.LookupEnv(key); ok {
		if valueInt, err := strconv.Atoi(value); err == nil {
			return valueInt
		}
	}
	return defaultValue
}

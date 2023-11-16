package utils

import (
	"log"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Version               string        `mapstructure:"VERSION"`
	Environment           string        `mapstructure:"ENVIRONMENT"`
	Port                  string        `mapstructure:"PORT"`
	ClientOrigin          string        `mapstructure:"CLIENT_ORIGIN"`
	DBHost                string        `mapstructure:"DB_HOST"`
	DBUser                string        `mapstructure:"DB_USER"`
	DBPassword            string        `mapstructure:"DB_PASSWORD"`
	DBName                string        `mapstructure:"DB_NAME"`
	DBPort                string        `mapstructure:"DB_PORT"`
	DBLogging             string        `mapstructure:"DB_LOGGING"`
	AccessTokenPrivateKey string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey  string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn  time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRES_IN"`
	AccessTokenMaxAge     int           `mapstructure:"ACCESS_TOKEN_MAX_AGE"`
}

var configInstance Config

func LoadConfig(path string) Config {
	err := godotenv.Load(path + "local.env")
	if err != nil {
		log.Println("Could not load .env.local file:", err)
		err = godotenv.Load(path + ".env")
		if err != nil {
			log.Println("No .env* files found, moving forward with host env")
		}
	}

	AccessTokenExpires, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_EXPIRES_IN"))
	if err != nil {
		log.Fatal("Error parsing ACCESS_TOKEN_EXPIRED_IN")
	}
	AccessTokenAge, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_MAX_AGE"))
	if err != nil {
		log.Fatal("Error parsing ACCESS_TOKEN_MAX_AGE")
	}

	config := Config{
		Version:               os.Getenv("VERSION"),
		Environment:           os.Getenv("ENVIRONMENT"),
		Port:                  os.Getenv("PORT"),
		ClientOrigin:          os.Getenv("CLIENT_ORIGIN"),
		DBHost:                os.Getenv("DB_HOST"),
		DBUser:                os.Getenv("DB_USER"),
		DBPassword:            os.Getenv("DB_PASSWORD"),
		DBName:                os.Getenv("DB_NAME"),
		DBPort:                os.Getenv("DB_PORT"),
		DBLogging:             os.Getenv("DB_LOGGING"),
		AccessTokenPrivateKey: os.Getenv("ACCESS_TOKEN_PRIVATE_KEY"),
		AccessTokenPublicKey:  os.Getenv("ACCESS_TOKEN_PUBLIC_KEY"),
		AccessTokenExpiresIn:  AccessTokenExpires,
		AccessTokenMaxAge:     AccessTokenAge,
	}

	testEnvsAreSet(config)

	SetConfig(config)
	return config
}

func SetConfig(config Config) {
	configInstance = config
}

func GetConfig() Config {
	return configInstance
}

func GetEnv(value string) string {
	return os.Getenv(value)
}

func testEnvsAreSet(config Config) {
	var unsetEnvs []string

	value := reflect.ValueOf(config)

	for i := 0; i < value.NumField(); i++ {
		if value.Field(i).Interface() == "" || (value.Field(i).Kind() == reflect.Int && value.Field(i).Int() == 0) {
			unsetEnvs = append(unsetEnvs, value.Type().Field(i).Tag.Get("mapstructure"))
		}
	}

	if len(unsetEnvs) > 0 {
		log.Fatalf("Shutting Down. The following environment variables are not set: %v", unsetEnvs)
	}
}

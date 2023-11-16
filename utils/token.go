package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type UserToken struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

func CreateToken(payload uuid.UUID) (string, error) {
	config := GetConfig()

	decodedPrivateKey, err := base64.StdEncoding.DecodeString(config.AccessTokenPrivateKey)

	if err != nil {
		return "", fmt.Errorf("could not decode key: %w", err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)

	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["sub"] = payload
	claims["exp"] = now.Add(config.AccessTokenExpiresIn).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)

	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

func ValidateToken(token string) (uuid.UUID, error) {
	config := GetConfig()
	decodedPublicKey, err := base64.StdEncoding.DecodeString(config.AccessTokenPublicKey)
	if err != nil {
		return uuid.Nil, fmt.Errorf("could not decode: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)

	if err != nil {
		return uuid.Nil, fmt.Errorf("validate: parse key: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return uuid.Nil, fmt.Errorf("validate: invalid token")
	}

	sub := claims["sub"].(uuid.UUID)

	return sub, nil
}

func DaysTokenValid(token string) int {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		log.Println("Invalid token format")
		return 0
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		log.Printf("Error decoding token payload: %v", err)
		return 0
	}

	var claims map[string]interface{}
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		log.Printf("Error unmarshalling token payload: %v", err)
		return 0
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		log.Println("Error: exp claim is not a float64")
		return 0
	}

	expTime := time.Unix(int64(exp), 0)
	return int(time.Until(expTime).Hours() / 24)
}

// func ExtractAndValidateUser(token string, config Config) (models.User, error) {
// 	sub, err := ValidateToken(token)
// 	if err != nil {
// 		return models.User{}, err
// 	}

// 	var user models.User
// 	DB := database.GetDatabase()
// 	if err := DB.Where("id = ?", sub).First(&user).Error; err != nil {
// 		return models.User{}, err
// 	}

// 	return user, nil
// }

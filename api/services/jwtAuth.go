package services

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

// JWT Service

type JWTService interface {
	GenerateJWT(UserId uint) (string, error)
	ValidateToken(encodedToken string) (*jwt.Token, error)
}

type authClaims struct {
	UserId uint `json:"user_id"`
	jwt.StandardClaims
}

type JWTServices struct {
	secretKey string
	issure    string
}

// Auth Service Constructor
func NewJWTAuthService() JWTService {
	return &JWTServices{
		secretKey: getSecretKey(),
		issure:    "Abinav",
	}
}

func getSecretKey() string {
	jwtSecretKey := os.Getenv("JWTSecretKey")
	if jwtSecretKey == "" {
		jwtSecretKey = "secretkey"
	}
	return jwtSecretKey
}

func (jwtServices *JWTServices) GenerateJWT(UserId uint) (string, error) {
	claims := &authClaims{
		UserId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			Issuer:    jwtServices.issure,
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//Encoded String
	//Signing Secret Key
	t, err := token.SignedString([]byte(jwtServices.secretKey))
	if err != nil {
		return "", fmt.Errorf("Error while generating token: %v", err.Error())
	}
	return t, nil
}

// Validate token from headers
func (jwtServices *JWTServices) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("Invalid Token %v", token.Header["alg"])
		}

		return []byte(jwtServices.secretKey), nil
	})
}

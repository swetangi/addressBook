package jwttoken

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

// var secretKey string = "secretKey123"
// var secretKeyByte = []byte(secretKey)

func GenerateJWT(email string) (string, error) {
	// token := jwt.New(jwt.SigningMethodHS256)
	// claims := token.Claims.(jwt.MapClaims)
	// claims["exp"] = time.Now().Add(10 * time.Minute)
	// claims["authorized"] = true
	// claims["user"] = "username"

	// or.....................

	// expirationTime := time.Now().Add(1 * time.Hour)

	// token := jwt.NewWithClaims(jwt.SigningMethodHS256,
	// 	&jwt.StandardClaims{
	// 		ExpiresAt: expirationTime.Unix(),
	// 		Subject:   email,
	// 	})

	//or.........................
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error in loading .env file")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})
	secretKey := os.Getenv("JWT_SECRET_KEY")
	secretKeyByte := []byte(secretKey)
	tokenString, err := token.SignedString(secretKeyByte)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenStr string) (jwt.Token, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error in loading .env file")
	}
	secretKey := os.Getenv("JWT_SECRET_KEY")
	secretKeyByte := []byte(secretKey)

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

		return secretKeyByte, nil
	})

	if err != nil {

		return jwt.Token{}, err
	}

	if !token.Valid {
		return jwt.Token{}, fmt.Errorf("invalid token")
	}

	return *token, nil
}

// func VerifyJWT(tokenString string) (bool, error) {
// 	token, err := jwt.Parse(request.Header) {

// 	})
// 	if err != nil {
// 		if err == jwt.ErrSignatureInvalid {
// 			return false, nil
// 		}
// 		return false, err
// 	}

// 	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
// 		if claims.ExpiresAt > time.Now().Unix() {
// 			return true, nil
// 		}
// 		return false, errors.New("token has expired")
// 	}
// 	return false, nil

// }

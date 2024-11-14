package middleware

import (
	jwttoken "addressBook/helpers/jwtToken"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Missing authorization header")
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwttoken.VerifyToken(tokenString)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Invalid token")
			r.Body.Close()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Failed to get token claims")
			return
		}

		// Get the email from the claims
		email, _ := claims["email"].(string)

		fmt.Println("Email extracted from token:", email)

		ctx := context.WithValue(r.Context(), "userEmail", email)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

	})
}

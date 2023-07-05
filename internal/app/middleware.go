package app

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (app *Application) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			next.ServeHTTP(w, r)
		} else {
			log.Println(r.Method)
			log.Println("Middleware exec")
			// Get cookie off request
			tokenString, err := r.Cookie("Authorization")
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			log.Println(tokenString.Value)

			// Decode/validate
			token, err := jwt.Parse(tokenString.Value, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("invalid signing method")
				}
				return []byte("SecretYouShouldHide"), nil
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				// Check the exp
				if float64(time.Now().Unix()) > claims["exp"].(float64) {
					http.Error(w, "token expired", http.StatusUnauthorized)
					return
				}

				// Find the user with token sub
				sub := fmt.Sprintf("%v", claims["sub"])
				if userLoader, err := app.DB.GetLoader(sub); err == nil {
					// Attach to request
					r.Header.Set("user", strconv.Itoa(userLoader.ID))
					r.Header.Set("role", "loader")

					// Continue
					next.ServeHTTP(w, r)
				} else if userCustomer, err := app.DB.GetCustomer(sub); err == nil {
					// Attach to request
					r.Header.Set("user", strconv.Itoa(userCustomer.ID))
					r.Header.Set("role", "customer")
					// Continue
					next.ServeHTTP(w, r)
				} else {
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
			} else {
				http.Error(w, "token invalid", http.StatusUnauthorized)
				return
			}
			log.Println("Middleware exec again")
		}
	})
}

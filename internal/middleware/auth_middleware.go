package middleware

import (
	"context"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"uzinfocom-todo/internal/models"
	"uzinfocom-todo/internal/models/response_objects"
)

func AuthJwtMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var responseObject *response_objects.ResponseObject
			tokenString := r.Header.Get("Authorization")

			if tokenString == "" {
				w.WriteHeader(http.StatusUnauthorized)
				responseObject = response_objects.NewResponseObject(false, "invalid token", nil)
				jsonResponse, _ := json.Marshal(responseObject)
				w.Write(jsonResponse)
				return
			}

			splitToken := strings.Split(tokenString, "Bearer ")
			reqToken := splitToken[1]

			token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
				return []byte(secretKey), nil
			})

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				responseObject = response_objects.NewResponseObject(false, "invalid token", nil)
				jsonResponse, _ := json.Marshal(responseObject)
				w.Write(jsonResponse)
				return
			}

			if !token.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				responseObject = response_objects.NewResponseObject(false, "invalid token", nil)
				jsonResponse, _ := json.Marshal(responseObject)
				w.Write(jsonResponse)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				responseObject = response_objects.NewResponseObject(false, "invalid token", nil)
				jsonResponse, _ := json.Marshal(responseObject)
				w.Write(jsonResponse)
				return
			}
			userId, _ := uuid.Parse(claims["user_id"].(string))
			name, _ := claims["name"].(string)
			phoneNumber, _ := claims["phone_number"].(string)
			user := models.User{
				UserId:      userId,
				Name:        name,
				PhoneNumber: phoneNumber,
			}
			ctx := context.WithValue(r.Context(), "user", user)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

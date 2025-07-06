package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/heinswanhtet/blogora-api/constants"
	store "github.com/heinswanhtet/blogora-api/stores"
	"github.com/heinswanhtet/blogora-api/types"
	"github.com/heinswanhtet/blogora-api/utils"
)

// func AuthenticateToken(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Println("checked authentication")

// 		next.ServeHTTP(w, r)
// 	})
// }

// type contextKey string
// const UserKey contextKey = "userID"
// const EmailKey contextKey = "email"

func AuthenticateToken(s *store.Store) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("checked authentication")
			token, err := utils.RetrieveBearerToken(r)
			if err != nil {
				log.Printf("failed to validate token: %v", err)
				permissionDenied(w)
				return
			}

			jwtToken, err := utils.ValidateJWT(token)
			if err != nil {
				log.Printf("failed to validate token: %v", err)
				permissionDenied(w)
				return
			}

			if !jwtToken.Valid {
				log.Println("invalid token")
				permissionDenied(w)
				return
			}

			claims := jwtToken.Claims.(jwt.MapClaims)
			userID := claims["userID"].(string)
			// fmt.Println(userID)

			u, err := s.GetAuthor(r.Context(), userID)
			if err != nil {
				log.Printf("failed to get user by id: %v", err)
				permissionDenied(w)
				return
			}
			// fmt.Println(*u.ID)

			// Add the user to the context
			ctx := r.Context()
			// ctx = context.WithValue(ctx, "userId", *u.ID)
			// ctx = context.WithValue(ctx, "email", *u.Email)
			ctx = context.WithValue(ctx, constants.ContextData, &types.ContextData{UserId: *u.ID, Email: *u.Email})
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied").Error())
}

// func GetUserIDFromContext(ctx context.Context) string {
// 	userID := ctx.Value(UserKey).(string)

// 	return userID
// }

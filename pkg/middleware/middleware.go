package middleware

import (
	"RestAPI/pkg/utils"
	"github.com/labstack/echo"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

func JWTMiddleware(secretKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
				return utils.JSONResponse(c, http.StatusUnauthorized, "error", "Missing or invalid token")
			}

			tokenString := authHeader[7:]
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrInvalidKey
				}
				return []byte(secretKey), nil
			})
			if err != nil || !token.Valid {
				return utils.JSONResponse(c, http.StatusUnauthorized, "error", "Invalid token")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return utils.JSONResponse(c, http.StatusUnauthorized, "error", "Invalid token claims")
			}

			// Сохраняем user_id в контексте, для использования в обработчиках
			c.Set("user_id", claims["user_id"])
			return next(c)
		}
	}
}

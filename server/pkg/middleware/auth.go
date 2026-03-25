package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"ims-intro/pkg/domain"
	"net/http"
	"os"
	"strings"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "missing or invalid token"})
		}

		tokenString := cookie.Value
		jwtKey := os.Getenv("JWT_KEY")

		claims := &domain.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid token"})
		}

		c.Set("user", claims)
		c.Set("role", strings.ToLower(claims.Role))
		c.Set("user_id", claims.UserID)
		return next(c)
	}
}

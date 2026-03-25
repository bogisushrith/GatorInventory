package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func Authorize(allowedRoles []string) echo.MiddlewareFunc {
	allowed := map[string]bool{}
	for _, role := range allowedRoles {
		allowed[strings.ToLower(role)] = true
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			roleValue := c.Get("role")
			role, ok := roleValue.(string)
			if !ok || role == "" {
				return c.JSON(http.StatusForbidden, map[string]string{"message": "forbidden"})
			}

			if !allowed[strings.ToLower(role)] {
				return c.JSON(http.StatusForbidden, map[string]string{"message": "forbidden"})
			}

			return next(c)
		}
	}

}

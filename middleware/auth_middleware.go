package middleware

import (
	"net/http"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/golang-jwt/jwt/v5"
)

// JWTMiddleware returns an Echo JWT middleware using the configured secret.
func JWTMiddleware() echo.MiddlewareFunc {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "spotsync-secret-key"
	}
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(secret),
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"message": "Unauthorized: " + err.Error(),
				"errors":  nil,
			})
		},
	})
}

// RequireRole returns middleware that enforces a specific role from the JWT.
func RequireRole(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			userRole, ok := claims["role"].(string)
			if !ok || userRole != role {
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"success": false,
					"message": "Forbidden: insufficient permissions",
					"errors":  nil,
				})
			}
			return next(c)
		}
	}
}

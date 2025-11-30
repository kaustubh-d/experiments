package app

import (
	"os"

	"github.com/labstack/echo/v4"
)

// ApiAuth is a simple bearer token auth middleware.
// Expects AUTH_TOKEN env var to be set.
func ApiAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")

		if auth == "" {
			return echo.NewHTTPError(401, "missing Authorization header")
		}
		const prefix = "Bearer "
		if len(auth) <= len(prefix) || auth[:len(prefix)] != prefix {
			return echo.NewHTTPError(401, "invalid Authorization header")
		}
		token := auth[len(prefix):]
		expected := os.Getenv("AUTH_TOKEN")

		if expected == "" {
			return echo.NewHTTPError(401, "server not configured for auth")
		}
		if token != expected {
			return echo.NewHTTPError(401, "invalid token")
		}

		return next(c)
	}
}

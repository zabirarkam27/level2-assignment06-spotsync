package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zabirarkam27/level2-assignment06-spotsync/dto"
)

func RequireRole(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole, ok := c.Get("role").(string)
			if !ok || userRole != role {
				return c.JSON(http.StatusForbidden, dto.APIResponse{
					Success: false,
					Message: "You don't have permission to perform this action",
				})
			}
			return next(c)
		}
	}
}

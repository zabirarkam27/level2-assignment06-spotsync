package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zabirarkam27/level2-assignment06-spotsync/dto"
	"github.com/zabirarkam27/level2-assignment06-spotsync/service"
)

func ok(c echo.Context, status int, message string, data interface{}) error {
	return c.JSON(status, dto.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func fail(c echo.Context, status int, message string, details interface{}) error {
	return c.JSON(status, dto.APIResponse{
		Success: false,
		Message: message,
		Errors:  details,
	})
}

func parseID(c echo.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		return 0, errors.New("invalid id")
	}
	return uint(id), nil
}

func statusFromServiceError(err error) (int, string) {
	switch {
	case errors.Is(err, service.ErrNotFound):
		return http.StatusNotFound, "Resource not found"
	case errors.Is(err, service.ErrForbidden):
		return http.StatusForbidden, "You don't have permission to perform this action"
	case errors.Is(err, service.ErrZoneFull):
		return http.StatusConflict, "Parking zone is full. No spots available."
	case errors.Is(err, service.ErrInvalidCredentials):
		return http.StatusUnauthorized, "Invalid email or password"
	case errors.Is(err, service.ErrEmailExists):
		return http.StatusBadRequest, "Email already registered"
	default:
		return http.StatusInternalServerError, "Internal server error"
	}
}

func userIDFromContext(c echo.Context) uint {
	userID, _ := c.Get("userID").(uint)
	return userID
}

func roleFromContext(c echo.Context) string {
	role, _ := c.Get("role").(string)
	return role
}

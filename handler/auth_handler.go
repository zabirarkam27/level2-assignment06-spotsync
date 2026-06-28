package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zabirarkam27/level2-assignment06-spotsync/dto"
	"github.com/zabirarkam27/level2-assignment06-spotsync/service"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return fail(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	if err := c.Validate(req); err != nil {
		return fail(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	user, err := h.service.Register(req)
	if err != nil {
		status, message := statusFromServiceError(err)
		return fail(c, status, message, nil)
	}

	return ok(c, http.StatusCreated, "User registered successfully", user)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return fail(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	if err := c.Validate(req); err != nil {
		return fail(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	response, err := h.service.Login(req)
	if err != nil {
		status, message := statusFromServiceError(err)
		return fail(c, status, message, nil)
	}

	return ok(c, http.StatusOK, "Login successful", response)
}

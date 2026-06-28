package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zabirarkam27/level2-assignment06-spotsync/dto"
	"github.com/zabirarkam27/level2-assignment06-spotsync/service"
)

type ReservationHandler struct {
	service service.ReservationService
}

func NewReservationHandler(service service.ReservationService) *ReservationHandler {
	return &ReservationHandler{service: service}
}

func (h *ReservationHandler) Create(c echo.Context) error {
	var req dto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return fail(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	if err := c.Validate(req); err != nil {
		return fail(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	reservation, err := h.service.Create(userIDFromContext(c), req)
	if err != nil {
		status, message := statusFromServiceError(err)
		return fail(c, status, message, nil)
	}

	return ok(c, http.StatusCreated, "Reservation confirmed successfully", reservation)
}

func (h *ReservationHandler) GetMine(c echo.Context) error {
	reservations, err := h.service.GetMine(userIDFromContext(c))
	if err != nil {
		status, message := statusFromServiceError(err)
		return fail(c, status, message, nil)
	}

	return ok(c, http.StatusOK, "My reservations retrieved successfully", reservations)
}

func (h *ReservationHandler) Cancel(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return fail(c, http.StatusBadRequest, "Invalid reservation id", nil)
	}

	if err := h.service.Cancel(userIDFromContext(c), id, roleFromContext(c)); err != nil {
		status, message := statusFromServiceError(err)
		return fail(c, status, message, nil)
	}

	return ok(c, http.StatusOK, "Reservation cancelled successfully", nil)
}

func (h *ReservationHandler) GetAll(c echo.Context) error {
	reservations, err := h.service.GetAll()
	if err != nil {
		status, message := statusFromServiceError(err)
		return fail(c, status, message, nil)
	}

	return ok(c, http.StatusOK, "All reservations retrieved successfully", reservations)
}

package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zabirarkam27/level2-assignment06-spotsync/dto"
	"github.com/zabirarkam27/level2-assignment06-spotsync/service"
)

type ZoneHandler struct {
	service service.ZoneService
}

func NewZoneHandler(service service.ZoneService) *ZoneHandler {
	return &ZoneHandler{service: service}
}

func (h *ZoneHandler) Create(c echo.Context) error {
	var req dto.CreateZoneRequest
	if err := c.Bind(&req); err != nil {
		return fail(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	if err := c.Validate(req); err != nil {
		return fail(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	zone, err := h.service.Create(req)
	if err != nil {
		status, message := statusFromServiceError(err)
		return fail(c, status, message, nil)
	}

	return ok(c, http.StatusCreated, "Parking zone created successfully", zone)
}

func (h *ZoneHandler) GetAll(c echo.Context) error {
	zones, err := h.service.GetAll()
	if err != nil {
		status, message := statusFromServiceError(err)
		return fail(c, status, message, nil)
	}

	return ok(c, http.StatusOK, "Parking zones retrieved successfully", zones)
}

func (h *ZoneHandler) GetOne(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return fail(c, http.StatusBadRequest, "Invalid zone id", nil)
	}

	zone, err := h.service.GetByID(id)
	if err != nil {
		status, message := statusFromServiceError(err)
		return fail(c, status, message, nil)
	}

	return ok(c, http.StatusOK, "Parking zone retrieved successfully", zone)
}

func (h *ZoneHandler) Update(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return fail(c, http.StatusBadRequest, "Invalid zone id", nil)
	}

	var req dto.UpdateZoneRequest
	if err := c.Bind(&req); err != nil {
		return fail(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	if err := c.Validate(req); err != nil {
		return fail(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	zone, err := h.service.Update(id, req)
	if err != nil {
		status, message := statusFromServiceError(err)
		return fail(c, status, message, nil)
	}

	return ok(c, http.StatusOK, "Parking zone updated successfully", zone)
}

func (h *ZoneHandler) Delete(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return fail(c, http.StatusBadRequest, "Invalid zone id", nil)
	}

	if err := h.service.Delete(id); err != nil {
		status, message := statusFromServiceError(err)
		return fail(c, status, message, nil)
	}

	return ok(c, http.StatusOK, "Parking zone deleted successfully", nil)
}

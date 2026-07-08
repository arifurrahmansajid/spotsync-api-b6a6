package handler

import (
	"net/http"
	"spotsync/dto"
	"spotsync/service"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ZoneHandler struct {
	zoneSvc service.ZoneService
}

func NewZoneHandler(zoneSvc service.ZoneService) *ZoneHandler {
	return &ZoneHandler{zoneSvc: zoneSvc}
}

func (h *ZoneHandler) CreateZone(c echo.Context) error {
	// Admin only — enforced in middleware
	var req dto.CreateZoneRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid request body", err.Error()))
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Validation failed", err.Error()))
	}

	zone, err := h.zoneSvc.CreateZone(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("Failed to create parking zone", err.Error()))
	}

	return c.JSON(http.StatusCreated, successResponse("Parking zone created successfully", zone))
}

func (h *ZoneHandler) GetAllZones(c echo.Context) error {
	zones, err := h.zoneSvc.GetAllZones()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("Failed to retrieve parking zones", err.Error()))
	}
	return c.JSON(http.StatusOK, successResponse("Parking zones retrieved successfully", zones))
}

func (h *ZoneHandler) GetZoneByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid zone ID", err.Error()))
	}

	zone, err := h.zoneSvc.GetZoneByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, errorResponse("Parking zone not found", nil))
		}
		return c.JSON(http.StatusInternalServerError, errorResponse("Failed to retrieve parking zone", err.Error()))
	}

	return c.JSON(http.StatusOK, successResponse("Parking zone retrieved successfully", zone))
}

func (h *ZoneHandler) UpdateZone(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid zone ID", err.Error()))
	}

	var req dto.UpdateZoneRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid request body", err.Error()))
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Validation failed", err.Error()))
	}

	zone, err := h.zoneSvc.UpdateZone(uint(id), req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, errorResponse("Parking zone not found", nil))
		}
		return c.JSON(http.StatusInternalServerError, errorResponse("Failed to update parking zone", err.Error()))
	}

	return c.JSON(http.StatusOK, successResponse("Parking zone updated successfully", zone))
}

func (h *ZoneHandler) DeleteZone(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid zone ID", err.Error()))
	}

	if err := h.zoneSvc.DeleteZone(uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, errorResponse("Parking zone not found", nil))
		}
		return c.JSON(http.StatusInternalServerError, errorResponse("Failed to delete parking zone", err.Error()))
	}

	return c.JSON(http.StatusOK, successResponse("Parking zone deleted successfully", nil))
}

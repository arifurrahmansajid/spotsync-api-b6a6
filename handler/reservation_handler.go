package handler

import (
	"errors"
	"net/http"
	"spotsync/dto"
	"spotsync/service"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ReservationHandler struct {
	reservationSvc service.ReservationService
}

func NewReservationHandler(reservationSvc service.ReservationService) *ReservationHandler {
	return &ReservationHandler{reservationSvc: reservationSvc}
}

// extractClaims reads user ID and role from the Echo JWT context.
func extractClaims(c echo.Context) (uint, string) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := uint(claims["id"].(float64))
	role := claims["role"].(string)
	return id, role
}

func (h *ReservationHandler) CreateReservation(c echo.Context) error {
	userID, _ := extractClaims(c)

	var req dto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid request body", err.Error()))
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Validation failed", err.Error()))
	}

	reservation, err := h.reservationSvc.CreateReservation(userID, req)
	if err != nil {
		if errors.Is(err, service.ErrZoneFull) {
			return c.JSON(http.StatusConflict, errorResponse("Reservation failed", "Parking zone is full"))
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, errorResponse("Reservation failed", "Parking zone not found"))
		}
		return c.JSON(http.StatusInternalServerError, errorResponse("Reservation failed", err.Error()))
	}

	return c.JSON(http.StatusCreated, successResponse("Reservation confirmed successfully", reservation))
}

func (h *ReservationHandler) GetMyReservations(c echo.Context) error {
	userID, _ := extractClaims(c)

	reservations, err := h.reservationSvc.GetMyReservations(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("Failed to retrieve reservations", err.Error()))
	}

	return c.JSON(http.StatusOK, successResponse("My reservations retrieved successfully", reservations))
}

func (h *ReservationHandler) CancelReservation(c echo.Context) error {
	userID, role := extractClaims(c)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid reservation ID", err.Error()))
	}

	if err := h.reservationSvc.CancelReservation(uint(id), userID, role); err != nil {
		if err.Error() == "forbidden" {
			return c.JSON(http.StatusForbidden, errorResponse("Access denied", "You can only cancel your own reservations"))
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, errorResponse("Reservation not found", nil))
		}
		if err.Error() == "reservation already cancelled" {
			return c.JSON(http.StatusBadRequest, errorResponse("Cannot cancel", err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, errorResponse("Failed to cancel reservation", err.Error()))
	}

	return c.JSON(http.StatusOK, successResponse("Reservation cancelled successfully", nil))
}

func (h *ReservationHandler) GetAllReservations(c echo.Context) error {
	// Admin only — enforced in middleware
	reservations, err := h.reservationSvc.GetAllReservations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("Failed to retrieve reservations", err.Error()))
	}

	return c.JSON(http.StatusOK, successResponse("All reservations retrieved successfully", reservations))
}

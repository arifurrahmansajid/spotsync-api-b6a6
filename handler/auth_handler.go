package handler

import (
	"net/http"
	"spotsync/dto"
	"spotsync/service"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authSvc service.AuthService
}

func NewAuthHandler(authSvc service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid request body", err.Error()))
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Validation failed", err.Error()))
	}

	user, err := h.authSvc.Register(req)
	if err != nil {
		if err.Error() == "email already registered" {
			return c.JSON(http.StatusBadRequest, errorResponse("Registration failed", err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, errorResponse("Registration failed", err.Error()))
	}

	return c.JSON(http.StatusCreated, successResponse("User registered successfully", user))
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid request body", err.Error()))
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Validation failed", err.Error()))
	}

	resp, err := h.authSvc.Login(req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errorResponse("Login failed", err.Error()))
	}

	return c.JSON(http.StatusOK, successResponse("Login successful", resp))
}

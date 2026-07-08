package main

import (
	"log"
	"os"

	"spotsync/config"
	"spotsync/handler"
	appMiddleware "spotsync/middleware"
	"spotsync/repository"
	"spotsync/service"
	appValidator "spotsync/validator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// ── 1. Connect & Migrate Database ────────────────────────────────────────
	db := config.ConnectDB()

	// ── 2. Dependency Injection (Repo → Service → Handler) ──────────────────
	// Repositories
	userRepo := repository.NewUserRepository(db)
	zoneRepo := repository.NewZoneRepository(db)
	reservationRepo := repository.NewReservationRepository(db)

	// Services
	authSvc := service.NewAuthService(userRepo)
	zoneSvc := service.NewZoneService(zoneRepo)
	reservationSvc := service.NewReservationService(reservationRepo, zoneRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authSvc)
	zoneHandler := handler.NewZoneHandler(zoneSvc)
	reservationHandler := handler.NewReservationHandler(reservationSvc)

	// ── 3. Echo Setup ─────────────────────────────────────────────────────────
	e := echo.New()
	e.Validator = appValidator.New()

	// Global middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// ── 4. Routes ─────────────────────────────────────────────────────────────
	api := e.Group("/api/v1")

	// Auth (public)
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	// Parking Zones
	zones := api.Group("/zones")
	zones.GET("", zoneHandler.GetAllZones)    // public
	zones.GET("/:id", zoneHandler.GetZoneByID) // public
	// Admin-only zone management
	zones.POST("", zoneHandler.CreateZone,
		appMiddleware.JWTMiddleware(),
		appMiddleware.RequireRole("admin"),
	)
	zones.PUT("/:id", zoneHandler.UpdateZone,
		appMiddleware.JWTMiddleware(),
		appMiddleware.RequireRole("admin"),
	)
	zones.DELETE("/:id", zoneHandler.DeleteZone,
		appMiddleware.JWTMiddleware(),
		appMiddleware.RequireRole("admin"),
	)

	// Reservations
	reservations := api.Group("/reservations", appMiddleware.JWTMiddleware())
	reservations.POST("", reservationHandler.CreateReservation)
	reservations.GET("/my-reservations", reservationHandler.GetMyReservations)
	reservations.DELETE("/:id", reservationHandler.CancelReservation)
	// Admin-only
	reservations.GET("", reservationHandler.GetAllReservations,
		appMiddleware.RequireRole("admin"),
	)

	// ── 5. Start Server ───────────────────────────────────────────────────────
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚗 SpotSync server starting on port %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

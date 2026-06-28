package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/zabirarkam27/level2-assignment06-spotsync/config"
	"github.com/zabirarkam27/level2-assignment06-spotsync/dto"
	"github.com/zabirarkam27/level2-assignment06-spotsync/handler"
	appmiddleware "github.com/zabirarkam27/level2-assignment06-spotsync/middleware"
	"github.com/zabirarkam27/level2-assignment06-spotsync/models"
	"github.com/zabirarkam27/level2-assignment06-spotsync/repository"
	"github.com/zabirarkam27/level2-assignment06-spotsync/service"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func corsAllowedOrigins() []string {
	origins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if origins == "" {
		return []string{"*"}
	}

	allowedOrigins := make([]string, 0)
	for _, origin := range strings.Split(origins, ",") {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			allowedOrigins = append(allowedOrigins, origin)
		}
	}

	if len(allowedOrigins) == 0 {
		return []string{"*"}
	}
	return allowedOrigins
}

func main() {
	_ = godotenv.Load()

	if len(os.Getenv("JWT_SECRET")) < 32 {
		log.Fatal("JWT_SECRET must be set and at least 32 characters long")
	}

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	migrationCtx, cancelMigration := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelMigration()
	if err := db.WithContext(migrationCtx).AutoMigrate(&models.User{}, &models.ParkingZone{}, &models.Reservation{}); err != nil {
		log.Fatal("failed to run migrations: ", err)
	}

	e := echo.New()
	e.HideBanner = true
	e.Validator = &CustomValidator{validator: validator.New()}
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}
		code := http.StatusInternalServerError
		message := "Internal server error"
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = http.StatusText(he.Code)
		}
		_ = c.JSON(code, dto.APIResponse{Success: false, Message: message})
	}

	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: corsAllowedOrigins(),
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAuthorization},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
	}))

	userRepo := repository.NewUserRepository(db)
	zoneRepo := repository.NewZoneRepository(db)
	reservationRepo := repository.NewReservationRepository(db)

	authService := service.NewAuthService(userRepo)
	zoneService := service.NewZoneService(zoneRepo)
	reservationService := service.NewReservationService(reservationRepo, zoneRepo)

	authHandler := handler.NewAuthHandler(authService)
	zoneHandler := handler.NewZoneHandler(zoneService)
	reservationHandler := handler.NewReservationHandler(reservationService)

	api := e.Group("/api/v1")
	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, dto.APIResponse{Success: true, Message: "SpotSync API is healthy"})
	})

	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	api.GET("/zones", zoneHandler.GetAll)
	api.GET("/zones/:id", zoneHandler.GetOne)

	protected := api.Group("")
	protected.Use(appmiddleware.JWTMiddleware)
	protected.POST("/reservations", reservationHandler.Create)
	protected.GET("/reservations/my-reservations", reservationHandler.GetMine)
	protected.DELETE("/reservations/:id", reservationHandler.Cancel)

	admin := api.Group("")
	admin.Use(appmiddleware.JWTMiddleware)
	admin.Use(appmiddleware.RequireRole("admin"))
	admin.POST("/zones", zoneHandler.Create)
	admin.PUT("/zones/:id", zoneHandler.Update)
	admin.DELETE("/zones/:id", zoneHandler.Delete)
	admin.GET("/reservations", reservationHandler.GetAll)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(e.Start(":" + port))
}

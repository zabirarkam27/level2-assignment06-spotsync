package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/zabirarkam27/level2-assignment06-spotsync/dto"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: "Missing or invalid authorization token",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &dto.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: "Invalid or expired token",
			})
		}

		claims, ok := token.Claims.(*dto.JWTClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: "Invalid token claims",
			})
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		return next(c)
	}
}

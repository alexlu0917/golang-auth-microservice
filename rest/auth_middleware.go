package rest

import (
	"microauth/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewAuthMiddleware(s domain.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			id := domain.AuthTokenID(ctx.Request().Header.Get("Authorization"))
			if err := s.Validate(ctx.Request().Context(), id); err != nil {
				return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "invaild token"})
			}

			return next(ctx)
		}
	}
}

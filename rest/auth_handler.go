package rest

import (
	"microauth/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	service domain.AuthService
}

func NewAuthHandler(s domain.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h AuthHandler) HandleLogin(c echo.Context) error {
	var r domain.AuthParams

	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	t, err := h.service.Authenticate(c.Request().Context(), r)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": t})
}

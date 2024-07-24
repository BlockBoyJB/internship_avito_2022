package v1

import (
	_ "avito_intership/docs"
	"avito_intership/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

func NewRouter(h *echo.Echo, services *service.Services) {
	h.Use(middleware.Recover())
	h.GET("/ping", ping)
	h.GET("/swagger/*", echoSwagger.WrapHandler)

	auth := authMiddleware{auth: services.Auth}
	h.GET("/token", auth.getToken)

	v1 := h.Group("/api/v1", auth.authHandler)
	newAccountRouter(v1.Group("/accounts"), services.Account)
	newReservationRouter(v1.Group("/reservations"), services.Reservation)
	newOperationRouter(v1.Group("/operations"), services.Operation)
}

func ping(c echo.Context) error {
	return c.NoContent(200)
}

// Вообще в микросервисе этого не должно быть, но так как я разрабатываю его изолированно, то можно
func (h *authMiddleware) getToken(c echo.Context) error {
	type response struct {
		Token string `json:"token"`
	}
	token, err := h.auth.CreateToken()
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
		return err
	}
	return c.JSON(http.StatusOK, response{Token: token})
}

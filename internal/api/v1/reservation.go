package v1

import (
	"avito_intership/internal/service"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type reservationRouter struct {
	reservation service.Reservation
}

func newReservationRouter(g *echo.Group, reservation service.Reservation) {
	r := &reservationRouter{reservation: reservation}

	g.POST("/create", r.create)
	g.DELETE("/cancel", r.cancel)
	g.POST("/revenue", r.revenue)
}

type reservationCreateInput struct {
	UserId    int     `json:"user_id" validate:"required"`
	ProductId int     `json:"product_id" validate:"required"`
	OrderId   int     `json:"order_id" validate:"required"`
	Amount    float64 `json:"amount" validate:"amount,required"`
}

type reservationResponse struct {
	ReservationId int `json:"reservation_id"`
}

//	@Summary		create reservation
//	@Description	Create product amount reservation
//	@Tags			reservation
//	@Accept			json
//	@Produce		json
//	@Param			input	body		reservationCreateInput	true	"input"
//	@Success		200		{object}	reservationResponse
//	@Failure		400		{object}	echo.HTTPError
//	@Failure		500		{object}	echo.HTTPError
//	@Security		JWT
//	@Router			/api/v1/reservations/create [post]
func (r *reservationRouter) create(c echo.Context) error {
	var input reservationCreateInput

	if err := c.Bind(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, echo.ErrBadRequest)
		return err
	}
	if err := c.Validate(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return err
	}

	reservationId, err := r.reservation.CreateReservation(c.Request().Context(), service.ReservationInput{
		UserId:    input.UserId,
		ProductId: input.ProductId,
		OrderId:   input.OrderId,
		Amount:    input.Amount,
	})
	if err != nil {
		if !errors.Is(err, service.ErrReservationCannotCreate) {
			errorResponse(c, http.StatusBadRequest, err)
			return nil
		}
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return err
	}

	return c.JSON(http.StatusCreated, reservationResponse{ReservationId: reservationId})
}

type reservationCancelInput struct {
	ReservationId int `json:"reservation_id"`
}

//	@Summary		cancel reservation
//	@Description	cancel product amount reservation and return money to account
//	@Tags			reservation
//	@Accept			json
//	@Produce		json
//	@Param			input	body	reservationCancelInput	true	"input"
//	@Success		200
//	@Failure		400	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Security		JWT
//	@Router			/api/v1/reservations/cancel [delete]
func (r *reservationRouter) cancel(c echo.Context) error {
	var input reservationCancelInput

	if err := c.Bind(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, echo.ErrBadRequest)
		return err
	}
	if err := c.Validate(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return err
	}

	if err := r.reservation.CancelReservation(c.Request().Context(), input.ReservationId); err != nil {
		if errors.Is(err, service.ErrReservationNotFound) {
			errorResponse(c, http.StatusBadRequest, err)
			return nil
		}
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return err
	}

	return c.NoContent(http.StatusOK)
}

type reservationRevenueInput struct {
	ReservationId int `json:"reservation_id"`
}

//	@Summary		revenue reservation
//	@Description	confirm reservation
//	@Tags			reservation
//	@Accept			json
//	@Produce		json
//	@Param			input	body	reservationRevenueInput	true	"input"
//	@Success		200
//	@Failure		400	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Security		JWT
//	@Router			/api/v1/reservations/revenue [post]
func (r *reservationRouter) revenue(c echo.Context) error {
	var input reservationRevenueInput

	if err := c.Bind(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, echo.ErrBadRequest)
		return err
	}
	if err := c.Validate(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return err
	}

	if err := r.reservation.RevenueReservation(c.Request().Context(), input.ReservationId); err != nil {
		if errors.Is(err, service.ErrReservationNotFound) {
			errorResponse(c, http.StatusBadRequest, err)
			return nil
		}
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return err
	}

	return c.NoContent(http.StatusOK)
}

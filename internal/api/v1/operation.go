package v1

import (
	"avito_intership/internal/service"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type operationRouter struct {
	operation service.Operation
}

func newOperationRouter(g *echo.Group, operation service.Operation) {
	r := &operationRouter{operation: operation}

	g.GET("/history", r.history)
	g.GET("/report", r.report)
}

type operationHistoryInput struct {
	UserId int    `json:"user_id" validate:"required"`
	Sort   string `json:"sort"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

//	@Summary		Get history
//	@Description	Get account transactions history
//	@Tags			operation
//	@Accept			json
//	@Produce		json
//	@Param			input	body		operationHistoryInput	true	"input"
//	@Success		200		{array}		service.HistoryOutput
//	@Failure		400		{object}	echo.HTTPError
//	@Failure		500		{object}	echo.HTTPError
//	@Security		JWT
//	@Router			/api/v1/operations/history [get]
func (r *operationRouter) history(c echo.Context) error {
	var input operationHistoryInput

	if err := c.Bind(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, echo.ErrBadRequest)
		return err
	}
	if err := c.Validate(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return err
	}

	history, err := r.operation.GetHistory(c.Request().Context(), service.HistoryInput{
		UserId: input.UserId,
		Sort:   input.Sort,
		Offset: input.Offset,
		Limit:  input.Limit,
	})
	if err != nil {
		if errors.Is(err, service.ErrAccountNotFound) {
			errorResponse(c, http.StatusBadRequest, err)
			return nil
		}
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return err
	}

	return c.JSON(http.StatusOK, history)
}

type operationReportInput struct {
	Year  int `json:"year" validate:"required"`
	Month int `json:"month" validate:"required"`
}

//	@Summary		Get report
//	@Description	Get monthly report, ordered by products ids
//	@Tags			operation
//	@Accept			json
//	@Produce		json
//	@Param			input	body	operationReportInput	true	"input"
//	@Success		200
//	@Failure		400	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Security		JWT
//	@Router			/api/v1/operations/report [get]
func (r *operationRouter) report(c echo.Context) error {
	var input operationReportInput

	if err := c.Bind(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, echo.ErrBadRequest)
		return err
	}
	if err := c.Validate(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return err
	}

	report, err := r.operation.CreateReport(c.Request().Context(), input.Year, input.Month)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return err
	}

	return c.Blob(http.StatusOK, "text/csv", report)
}

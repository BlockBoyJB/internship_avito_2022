package v1

import (
	"avito_intership/internal/service"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type accountRouter struct {
	account service.Account
}

func newAccountRouter(g *echo.Group, account service.Account) {
	r := &accountRouter{account: account}

	g.POST("/create", r.create)
	g.GET("/balance", r.balance)
	g.PATCH("/deposit", r.deposit)
	g.PATCH("/withdraw", r.withdraw)
	g.POST("/transfer", r.transfer)
}

type accountCreateInput struct {
	UserId int `json:"user_id" validate:"required"`
}

// @Summary		Create account
// @Description	Create account
// @Tags			account
// @Accept			json
// @Produce		json
// @Param			input	body	accountCreateInput	true	"input"
// @Success		201
// @Failure		400	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Security		JWT
// @Router			/api/v1/accounts/create [post]
func (r *accountRouter) create(c echo.Context) error {
	var input accountCreateInput

	if err := c.Bind(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, echo.ErrBadRequest)
		return err
	}
	if err := c.Validate(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return err
	}

	if err := r.account.CreateAccount(c.Request().Context(), input.UserId); err != nil {
		if !errors.Is(err, service.ErrAccountCannotCreate) {
			errorResponse(c, http.StatusBadRequest, err)
			return nil
		}
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return err
	}
	return c.NoContent(http.StatusCreated)
}

type balanceResponse struct {
	Balance float64 `json:"balance"`
}

// @Summary		Get balance
// @Description	Get balance for account by id
// @Tags			account
// @Accept			json
// @Produce		json
// @Param			user_id	query		string	true	"user id"
// @Success		200		{object}	balanceResponse
// @Failure		400		{object}	echo.HTTPError
// @Failure		500		{object}	echo.HTTPError
// @Security		JWT
// @Router			/api/v1/accounts/balance [get]
func (r *accountRouter) balance(c echo.Context) error {
	q := c.QueryParam("user_id") // возможно, стоит указывать в теле запроса
	if len(q) == 0 {
		errorResponse(c, http.StatusBadRequest, echo.ErrBadRequest)
		return nil
	}
	userId, err := strconv.Atoi(q)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, echo.ErrBadRequest)
		return err
	}

	balance, err := r.account.GetBalance(c.Request().Context(), userId)
	if err != nil {
		if errors.Is(err, service.ErrAccountNotFound) {
			errorResponse(c, http.StatusBadRequest, err)
			return nil
		}
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return err
	}

	return c.JSON(http.StatusOK, balanceResponse{Balance: balance})
}

type accountDepositInput struct {
	UserId int     `json:"user_id" validate:"required"`
	Amount float64 `json:"amount" validate:"amount,required"`
}

// @Summary		Account deposit
// @Description	Deposit on account
// @Tags			account
// @Accept			json
// @Produce		json
// @Param			input	body	accountDepositInput	true	"input"
// @Success		200
// @Failure		400	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Security		JWT
// @Router			/api/v1/accounts/deposit [patch]
func (r *accountRouter) deposit(c echo.Context) error {
	var input accountDepositInput

	if err := c.Bind(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, echo.ErrBadRequest)
		return err
	}
	if err := c.Validate(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return err
	}

	err := r.account.Deposit(c.Request().Context(), service.DepositInput{
		UserId: input.UserId,
		Amount: input.Amount,
	})
	if err != nil {
		if errors.Is(err, service.ErrAccountNotFound) {
			errorResponse(c, http.StatusBadRequest, err)
			return nil
		}
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return err
	}
	return c.NoContent(http.StatusOK)
}

type accountWithdrawInput struct {
	UserId int     `json:"user_id" validate:"required"`
	Amount float64 `json:"amount" validate:"amount,required"`
}

// @Summary		Account withdraw
// @Description	Withdraw from account
// @Tags			account
// @Accept			json
// @Produce		json
// @Param			input	body	accountWithdrawInput	true	"input"
// @Success		200
// @Failure		400	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Security		JWT
// @Router			/api/v1/accounts/withdraw [patch]
func (r *accountRouter) withdraw(c echo.Context) error {
	var input accountWithdrawInput

	if err := c.Bind(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, echo.ErrBadRequest)
		return err
	}
	if err := c.Validate(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return err
	}

	err := r.account.Withdraw(c.Request().Context(), service.WithdrawInput{
		UserId: input.UserId,
		Amount: input.Amount,
	})
	if err != nil {
		if !errors.Is(err, service.ErrCannotUpdateBalance) {
			errorResponse(c, http.StatusBadRequest, err)
			return nil
		}
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return err
	}

	return c.NoContent(http.StatusOK)
}

type accountTransferInput struct {
	From   int     `json:"from" validate:"required"`
	To     int     `json:"to" validate:"required"`
	Amount float64 `json:"amount" validate:"amount,required"`
}

// @Summary		Account transfer
// @Description	Transfer from account to account
// @Tags			account
// @Accept			json
// @Produce		json
// @Param			input	body	accountTransferInput	true	"input"
// @Success		200
// @Failure		400	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Security		JWT
// @Router			/api/v1/accounts/transfer [post]
func (r *accountRouter) transfer(c echo.Context) error {
	var input accountTransferInput

	if err := c.Bind(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, echo.ErrBadRequest)
		return err
	}
	if err := c.Validate(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return err
	}

	err := r.account.Transfer(c.Request().Context(), service.TransferInput{
		From:   input.From,
		To:     input.To,
		Amount: input.Amount,
	})
	if err != nil {
		if !errors.Is(err, service.ErrCannotUpdateBalance) {
			errorResponse(c, http.StatusBadRequest, err)
			return nil
		}
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return err
	}

	return c.NoContent(http.StatusOK)
}

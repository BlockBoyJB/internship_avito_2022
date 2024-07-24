package service

import "errors"

var (
	ErrAccountAlreadyExists = errors.New("account already exists")
	ErrAccountCannotCreate  = errors.New("cannot create account")
	ErrAccountNotFound      = errors.New("account not found")

	ErrNotEnoughBalance    = errors.New("not enough balance on account")
	ErrCannotUpdateBalance = errors.New("cannot update account balance")

	ErrReservationCannotCreate = errors.New("cannot create reservation")
	ErrReservationNotFound     = errors.New("reservation not found")
)

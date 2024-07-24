package service

import (
	"avito_intership/internal/model/dbmodel"
	"avito_intership/internal/repo"
	"avito_intership/internal/repo/pgerrs"
	"avito_intership/pkg/broker"
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
)

const (
	reservationPrefixLog = "/service/reservation"
)

type reservationService struct {
	reservation repo.Reservation
	producer    broker.Producer
}

func newReservationService(reservation repo.Reservation, producer broker.Producer) *reservationService {
	return &reservationService{
		reservation: reservation,
		producer:    producer,
	}
}

func (s *reservationService) CreateReservation(ctx context.Context, input ReservationInput) (int, error) {
	reservationId, err := s.reservation.CreateReservation(ctx, dbmodel.Reservation{
		UserId:    input.UserId,
		ProductId: input.ProductId,
		OrderId:   input.OrderId,
		Amount:    input.Amount,
	})
	if err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			return 0, ErrAccountNotFound
		}
		if errors.Is(err, pgerrs.ErrNotEnoughBalance) {
			return 0, ErrNotEnoughBalance
		}
		log.Errorf("%s/CreateReservation error create reservation: %s", reservationPrefixLog, err)
		return 0, ErrReservationCannotCreate
	}

	err = pushMessage(s.producer, dbmodel.OperationReservation, message{
		UserId: input.UserId,
		Amount: input.Amount,
	})
	if err != nil {
		return 0, ErrReservationCannotCreate
	}
	return reservationId, nil
}

func (s *reservationService) CancelReservation(ctx context.Context, reservationId int) error {
	userId, amount, err := s.reservation.DeleteReservation(ctx, reservationId)
	if err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			return ErrReservationNotFound
		}
		log.Errorf("%s/CancelReservation error delete reservation: %s", reservationPrefixLog, err)
		return err
	}

	err = pushMessage(s.producer, dbmodel.OperationDereservation, message{
		UserId: userId,
		Amount: amount,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *reservationService) RevenueReservation(ctx context.Context, reservationId int) error {
	userId, amount, err := s.reservation.RevenueReservation(ctx, reservationId)
	if err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			return ErrReservationNotFound
		}
		log.Errorf("%s/RevenueReservation error refund recognition: %s", reservationPrefixLog, err)
		return err
	}

	err = pushMessage(s.producer, dbmodel.OperationRevenue, message{
		UserId: userId,
		Amount: amount,
	})
	if err != nil {
		return err
	}
	return nil
}

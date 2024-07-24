package pgdb

import (
	"avito_intership/internal/model/dbmodel"
	"avito_intership/internal/repo/pgerrs"
	"avito_intership/pkg/postgres"
	"avito_intership/pkg/redis"
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

const reservationPrefixLog = "/pgdb/reservation"

type ReservationRepo struct {
	*postgres.Postgres
	redis redis.Redis
}

func NewReservationRepo(pg *postgres.Postgres, redis redis.Redis) *ReservationRepo {
	return &ReservationRepo{
		Postgres: pg,
		redis:    redis,
	}
}

func (r *ReservationRepo) CreateReservation(ctx context.Context, reservation dbmodel.Reservation) (int, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		log.Errorf("%s/CreateReservation error init tx: %s", reservationPrefixLog, err)
		return 0, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	balance, err := getCacheBalance(ctx, r.redis, reservation.UserId)
	if err != nil {
		if !errors.Is(err, pgerrs.ErrNotFound) {
			return 0, err
		}
		balance, err = getBalanceTx(ctx, tx, r.Builder, reservation.UserId)
		if err != nil {
			return 0, err
		}
	}

	if balance < reservation.Amount {
		return 0, pgerrs.ErrNotEnoughBalance
	}

	sql, args, _ := r.Builder.
		Update("account").
		Set("balance", squirrel.Expr("balance - ?", reservation.Amount)).
		Where("user_id = ?", reservation.UserId).
		Suffix("returning balance").
		ToSql()

	if err = tx.QueryRow(ctx, sql, args...).Scan(&balance); err != nil {
		log.Errorf("%s/CreateReservation error update account balance: %s", reservationPrefixLog, err)
		return 0, err
	}

	if err = setCacheBalance(ctx, r.redis, reservation.UserId, balance); err != nil {
		return 0, err
	}

	var reservationId int
	sql, args, _ = r.Builder.
		Insert("reservation").
		Columns("user_id", "product_id", "order_id", "amount").
		Values(reservation.UserId, reservation.ProductId, reservation.OrderId, reservation.Amount).
		Suffix("returning id").
		ToSql()

	if err = tx.QueryRow(ctx, sql, args...).Scan(&reservationId); err != nil {
		log.Errorf("%s/CreateReservation error create reservation: %s", reservationPrefixLog, err)
		return 0, err
	}

	sql, args, _ = r.Builder.
		Insert("operation").
		Columns("user_id", "product_id", "order_id", "amount", "type").
		Values(reservation.UserId, reservation.ProductId, reservation.OrderId, reservation.Amount, dbmodel.OperationReservation).
		ToSql()
	if _, err = tx.Exec(ctx, sql, args...); err != nil {
		log.Errorf("%s/CreateReservation error create operation: %s", reservationPrefixLog, err)
		return 0, err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Errorf("%s/CreateReservation error commit: %s", reservationPrefixLog, err)
		return 0, err
	}
	return reservationId, nil
}

func (r *ReservationRepo) DeleteReservation(ctx context.Context, reservationId int) (int, float64, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		log.Errorf("%s/DeleteReservation error init tx: %s", reservationPrefixLog, err)
		return 0, 0, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	sql, args, _ := r.Builder.
		Delete("reservation").
		Where("id = ?", reservationId).
		Suffix("returning user_id, product_id, order_id, amount").
		ToSql()

	var reservation dbmodel.Reservation
	if err = tx.QueryRow(ctx, sql, args...).Scan(&reservation.UserId, &reservation.ProductId, &reservation.OrderId, &reservation.Amount); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, 0, pgerrs.ErrNotFound
		}
		log.Errorf("%s/DeleteReservation error delete reservation: %s", reservationPrefixLog, err)
		return 0, 0, err
	}

	sql, args, _ = r.Builder.
		Update("account").
		Set("balance", squirrel.Expr("balance + ?", reservation.Amount)).
		Where("user_id = ?", reservation.UserId).
		Suffix("returning balance").
		ToSql()

	var balance float64
	if err = tx.QueryRow(ctx, sql, args...).Scan(&balance); err != nil {
		log.Errorf("%s/DeleteReservation error update account balance: %s", reservationPrefixLog, err)
		return 0, 0, err
	}

	if err = setCacheBalance(ctx, r.redis, reservation.UserId, balance); err != nil {
		return 0, 0, err
	}

	sql, args, _ = r.Builder.
		Insert("operation").
		Columns("user_id", "product_id", "order_id", "amount", "type").
		Values(reservation.UserId, reservation.ProductId, reservation.OrderId, reservation.Amount, dbmodel.OperationDereservation).
		ToSql()
	if _, err = tx.Exec(ctx, sql, args...); err != nil {
		log.Errorf("%s/DeleteReservation error create operation: %s", reservationPrefixLog, err)
		return 0, 0, err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Errorf("%s/DeleteReservation error commit: %s", reservationPrefixLog, err)
		return 0, 0, err
	}
	return reservation.UserId, reservation.Amount, nil
}

func (r *ReservationRepo) RevenueReservation(ctx context.Context, reservationId int) (int, float64, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		log.Errorf("%s/RevenueReservation error init tx: %s", reservationPrefixLog, err)
		return 0, 0, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	sql, args, _ := r.Builder.
		Delete("reservation").
		Where("id = ?", reservationId).
		Suffix("returning user_id, product_id, order_id, amount").
		ToSql()

	var reservation dbmodel.Reservation
	if err = tx.QueryRow(ctx, sql, args...).Scan(&reservation.UserId, &reservation.ProductId, &reservation.OrderId, &reservation.Amount); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, 0, pgerrs.ErrNotFound
		}
		log.Errorf("%s/RevenueReservation error delete reservation: %s", reservationPrefixLog, err)
		return 0, 0, err
	}

	sql, args, _ = r.Builder.
		Insert("operation").
		Columns("user_id", "product_id", "order_id", "amount", "type").
		Values(reservation.UserId, reservation.ProductId, reservation.OrderId, reservation.Amount, dbmodel.OperationRevenue).
		ToSql()
	if _, err = tx.Exec(ctx, sql, args...); err != nil {
		log.Errorf("%s/RevenueReservation error create operation: %s", reservationPrefixLog, err)
		return 0, 0, err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Errorf("%s/RevenueReservation error commit: %s", reservationPrefixLog, err)
		return 0, 0, err
	}
	return reservation.UserId, reservation.Amount, nil
}

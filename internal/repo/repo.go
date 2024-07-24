package repo

import (
	"avito_intership/internal/model/dbmodel"
	"avito_intership/internal/repo/pgdb"
	"avito_intership/pkg/postgres"
	"avito_intership/pkg/redis"
	"context"
)

type Account interface {
	CreateAccount(ctx context.Context, userId int) error
	GetBalance(ctx context.Context, userId int) (float64, error)

	Deposit(ctx context.Context, userId int, amount float64) error
	Withdraw(ctx context.Context, userId int, amount float64) error
	Transfer(ctx context.Context, sendId, receiveId int, amount float64) error
}

type Reservation interface {
	CreateReservation(ctx context.Context, reservation dbmodel.Reservation) (int, error)
	DeleteReservation(ctx context.Context, reservationId int) (int, float64, error)
	RevenueReservation(ctx context.Context, reservationId int) (int, float64, error)
}

type Operation interface {
	GetHistory(ctx context.Context, userId int, sort string, offset, limit int) ([]dbmodel.Operation, error)
	GroupProductRevenue(ctx context.Context, year, month int) (map[int]float64, error)
}

type Repositories struct {
	Account
	Reservation
	Operation
}

func NewRepositories(pg *postgres.Postgres, redis redis.Redis) *Repositories {
	return &Repositories{
		Account:     pgdb.NewAccountRepo(pg, redis),
		Reservation: pgdb.NewReservationRepo(pg, redis),
		Operation:   pgdb.NewOperationRepo(pg),
	}
}

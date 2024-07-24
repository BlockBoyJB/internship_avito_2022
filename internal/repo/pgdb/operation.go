package pgdb

import (
	"avito_intership/internal/model/dbmodel"
	"avito_intership/internal/repo/pgerrs"
	"avito_intership/pkg/postgres"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

const operationPrefixLog = "/pgdb/operation"

type OperationRepo struct {
	*postgres.Postgres
}

func NewOperationRepo(pg *postgres.Postgres) *OperationRepo {
	return &OperationRepo{pg}
}

func (r *OperationRepo) GetHistory(ctx context.Context, userId int, sort string, offset, limit int) ([]dbmodel.Operation, error) {
	sql, args, _ := r.Builder.
		Select("*").
		From("operation").
		Where("user_id = ?", userId).
		OrderBy(sort).
		Offset(uint64(offset)).
		Limit(uint64(limit)).
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgerrs.ErrNotFound
		}
		log.Errorf("%s/GetHistory error get operations: %s", operationPrefixLog, err)
		return nil, err
	}
	defer rows.Close()

	var result []dbmodel.Operation
	for rows.Next() {
		var operation dbmodel.Operation

		err = rows.Scan(
			&operation.Id,
			&operation.UserId,
			&operation.ProductId,
			&operation.OrderId,
			&operation.Amount,
			&operation.Type,
			&operation.CreatedAt,
		)
		if err != nil {
			log.Errorf("%s/GetHistory error get operation: %s", operationPrefixLog, err)
			continue
		}
		result = append(result, operation)
	}
	return result, nil
}

func (r *OperationRepo) GroupProductRevenue(ctx context.Context, year, month int) (map[int]float64, error) {
	sql, args, _ := r.Builder.
		Select("product_id", "sum(amount)").
		From("operation").
		Where("type = ? and extract(year from operation.created_at) = ? and extract(month from operation.created_at) = ?", dbmodel.OperationRevenue, year, month).
		GroupBy("product_id").
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		log.Errorf("%s/GroupProductRevenue error get products: %s", operationPrefixLog, err)
		return nil, err
	}
	defer rows.Close()

	result := make(map[int]float64)

	for rows.Next() {
		var productId int
		var amount float64

		if err = rows.Scan(&productId, &amount); err != nil {
			log.Errorf("%s/GroupProductRevenue error get product: %s", operationPrefixLog, err)
			continue
		}
		result[productId] = amount
	}
	return result, nil
}

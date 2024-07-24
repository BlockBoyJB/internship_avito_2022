package pgdb

import (
	"avito_intership/internal/model/dbmodel"
	"avito_intership/internal/repo/pgerrs"
	"avito_intership/pkg/postgres"
	"avito_intership/pkg/redis"
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	accountPrefixLog = "/pgdb/account"
	defaultBalanceTL = time.Hour * 72
)

var key = func(id int) string { return fmt.Sprintf("balance:%d", id) }

type AccountRepo struct {
	*postgres.Postgres
	redis redis.Redis
}

func NewAccountRepo(pg *postgres.Postgres, redis redis.Redis) *AccountRepo {
	return &AccountRepo{
		Postgres: pg,
		redis:    redis,
	}
}

func (r *AccountRepo) CreateAccount(ctx context.Context, userId int) error {
	sql, args, _ := r.Builder.
		Insert("account").
		Columns("user_id").
		Values(userId).
		ToSql()

	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return pgerrs.ErrAlreadyExists
			}
		}
		log.Errorf("%s/CreateAccount error exec stmt: %s", accountPrefixLog, err)
		return err
	}
	return nil
}

// GetBalance смотрим сначала в кэш. Если нет, то идем в базу, там получаем. В конце пытаемся записать баланс в кэш
// Ошибка не хэндлится, потому что не критично, если не запишем
func (r *AccountRepo) GetBalance(ctx context.Context, userId int) (float64, error) {
	balance, err := getCacheBalance(ctx, r.redis, userId)
	if err == nil {
		return balance, nil
	}
	if err != nil {
		if !errors.Is(err, pgerrs.ErrNotFound) {
			return 0, err
		}
	}

	sql, args, _ := r.Builder.
		Select("balance").
		From("account").
		Where("user_id = ?", userId).
		ToSql()

	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(&balance); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, pgerrs.ErrNotFound
		}
		log.Errorf("%s/GetBalance error get balance: %s", accountPrefixLog, err)
		return 0, err
	}

	_ = setCacheBalance(ctx, r.redis, userId, balance)

	return balance, nil
}

// Получение баланса из кэша. Если не найдено, то возвращает ошибку ErrNotFound
func getCacheBalance(ctx context.Context, redis redis.Redis, userId int) (float64, error) {
	ok, err := redis.Exists(ctx, key(userId)).Result()
	if err != nil {
		log.Errorf("%s/getCacheBalance error check user balance exist: %s", accountPrefixLog, err)
		return 0, err
	}
	if ok == 0 {
		return 0, pgerrs.ErrNotFound
	}
	balance, err := redis.Get(ctx, key(userId)).Float64()
	if err != nil {
		log.Errorf("%s/getCacheBalance error get balance: %s", accountPrefixLog, err)
		return 0, err
	}
	return balance, nil
}

// Сохранение (обновление) баланса в кэш. Дефолтное время хранения - 3 дня.
func setCacheBalance(ctx context.Context, redis redis.Redis, userId int, amount float64) error {
	if err := redis.Set(ctx, key(userId), amount, defaultBalanceTL).Err(); err != nil {
		log.Errorf("%s/setCacheBalance error set balance to cache: %s", accountPrefixLog, err)
		return err
	}
	return nil
}

// Вынес в отдельную функцию получение баланса. Сделано чисто под транзакции, хотя даже там можно использовать обычный GetBalance (наверно)
func getBalanceTx(ctx context.Context, tx pgx.Tx, builder squirrel.StatementBuilderType, userId int) (float64, error) {
	var balance float64

	sql, args, _ := builder.
		Select("balance").
		From("account").
		Where("user_id = ?", userId).
		ToSql()

	if err := tx.QueryRow(ctx, sql, args...).Scan(&balance); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, pgerrs.ErrNotFound
		}
		log.Errorf("%s/getBalanceTx error get balance: %s", accountPrefixLog, err)
		return 0, err
	}
	return balance, nil
}

func (r *AccountRepo) Deposit(ctx context.Context, userId int, amount float64) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		log.Errorf("%s/Deposit error init tx: %s", accountPrefixLog, err)
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var balance float64
	sql, args, _ := r.Builder.
		Update("account").
		Set("balance", squirrel.Expr("balance + ?", amount)).
		Where("user_id = ?", userId).
		Suffix("returning balance").
		ToSql()

	if err = tx.QueryRow(ctx, sql, args...).Scan(&balance); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgerrs.ErrNotFound
		}
		log.Errorf("%s/Deposit error update account balance: %s", accountPrefixLog, err)
		return err
	}

	if err = setCacheBalance(ctx, r.redis, userId, balance); err != nil {
		return err
	}

	sql, args, _ = r.Builder.
		Insert("operation").
		Columns("user_id", "amount", "type").
		Values(userId, amount, dbmodel.OperationDeposit).
		ToSql()

	if _, err = tx.Exec(ctx, sql, args...); err != nil {
		log.Errorf("%s/Deposit error create operation: %s", accountPrefixLog, err)
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		log.Errorf("%s/Deposit error commit: %s", accountPrefixLog, err)
		return err
	}
	return nil
}

func (r *AccountRepo) Withdraw(ctx context.Context, userId int, amount float64) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		log.Errorf("%s/Withdraw error init tx: %s", accountPrefixLog, err)
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	balance, err := getCacheBalance(ctx, r.redis, userId)
	if err != nil {
		if !errors.Is(err, pgerrs.ErrNotFound) {
			return err
		}
		balance, err = getBalanceTx(ctx, tx, r.Builder, userId)
		if err != nil {
			return err
		}
	}

	if balance < amount {
		return pgerrs.ErrNotEnoughBalance
	}

	sql, args, _ := r.Builder.
		Update("account").
		Set("balance", squirrel.Expr("balance - ?", amount)).
		Where("user_id = ?", userId).
		Suffix("returning balance").
		ToSql()

	if err = tx.QueryRow(ctx, sql, args...).Scan(&balance); err != nil {
		log.Errorf("%s/Withdraw error update account balance: %s", accountPrefixLog, err)
		return err
	}

	if err = setCacheBalance(ctx, r.redis, userId, balance); err != nil {
		return err
	}

	sql, args, _ = r.Builder.
		Insert("operation").
		Columns("user_id", "amount", "type").
		Values(userId, amount, dbmodel.OperationWithdraw).
		ToSql()

	if _, err = tx.Exec(ctx, sql, args...); err != nil {
		log.Errorf("%s/Withdraw error create operation: %s", accountPrefixLog, err)
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		log.Errorf("%s/Withdraw error commit: %s", accountPrefixLog, err)
		return err
	}
	return nil
}

func (r *AccountRepo) Transfer(ctx context.Context, sendId, receiveId int, amount float64) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		log.Errorf("%s/Transfer error init tx: %s", accountPrefixLog, err)
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	balance, err := getCacheBalance(ctx, r.redis, sendId)
	if err != nil {
		if !errors.Is(err, pgerrs.ErrNotFound) {
			return err
		}
		balance, err = getBalanceTx(ctx, tx, r.Builder, sendId)
		if err != nil {
			return err
		}
	}

	if balance < amount {
		return pgerrs.ErrNotEnoughBalance
	}

	sql, args, _ := r.Builder.
		Update("account").
		Set("balance", squirrel.Expr("balance - ?", amount)).
		Where("user_id = ?", sendId).
		Suffix("returning balance").
		ToSql()

	if err = tx.QueryRow(ctx, sql, args...).Scan(&balance); err != nil {
		log.Errorf("%s/Transfer error update sender account balance: %s", accountPrefixLog, err)
		return err
	}

	if err = setCacheBalance(ctx, r.redis, sendId, balance); err != nil {
		return err
	}

	sql, args, _ = r.Builder.
		Update("account").
		Set("balance", squirrel.Expr("balance + ?", amount)).
		Where("user_id = ?", receiveId).
		Suffix("returning balance").
		ToSql()

	if err = tx.QueryRow(ctx, sql, args...).Scan(&balance); err != nil {
		log.Errorf("%s/Transfer error update receiver account balance: %s", accountPrefixLog, err)
		return err
	}

	// важно обновить (создать) новый баланс в кэше, потому что если до этого существовало какое-то значение,
	// то появится проблема несоответствия значений в основной бд и кэше
	if err = setCacheBalance(ctx, r.redis, receiveId, balance); err != nil {
		return err
	}

	sql, args, _ = r.Builder.
		Insert("operation").
		Columns("user_id", "amount", "type").
		Values(sendId, amount, dbmodel.OperationOutgoingTransfer).
		ToSql()

	if _, err = tx.Exec(ctx, sql, args...); err != nil {
		log.Errorf("%s/Transfer error create operation: %s", accountPrefixLog, err)
		return err
	}
	args[0], args[2] = receiveId, dbmodel.OperationIncomingTransfer // нужно записать туда и обратно
	if _, err = tx.Exec(ctx, sql, args...); err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23503" { // Ошибка несуществующего получателя появляется только в этом месте
				return pgerrs.ErrNotFound
			}
		}
		log.Errorf("%s/Transfer error create operation: %s", accountPrefixLog, err)
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Errorf("%s/Transfer error commit: %s", accountPrefixLog, err)
		return err
	}
	return nil
}

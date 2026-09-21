package transaction

import (
	"context"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type ctxKey string

var txKey ctxKey = "tx"

type Transactor struct {
	db *sqlx.DB
}

func NewTransactor(db *sqlx.DB) *Transactor {
	return &Transactor{db: db}
}

// TODO: добавить unit-тесты на pkg/ (в первую очередь на Transaction —
// сценарии commit/rollback/panic на моке sqlx.DB/sqlx.Tx)

// Transaction UoW (Unit of Work) - паттерн, который гарантирует атомарность операций,
// затрагивающих несколько Repository в рамках одной бизнес-транзакции:
// либо все SQL-изменения внутри блока применяются успешно (commit),
// либо откатываются целиком при любой ошибке (rollback). Repository,
// вызванные внутри одного UoW-блока, работают с одной и той же SQL-транзакцией через общий контекст,
// вместо того чтобы каждый открывал и коммитил свою отдельную транзакцию независимо.
func (t *Transactor) Transaction(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	//Выполняем fn() в рамках пришедшей
	//транзакции с контекста, не создавая новую
	_, ok := ExtractTx(ctx)
	if ok {
		return fn(ctx)
	}
	trx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		slog.Error("failed to start transaction", err)
		return err
	}
	defer func() {
		//Recovery - реализует проверку паники внутри SQL запроса,
		//обернутого в транзакцию, если паника поймана - производит откат
		if rec := recover(); rec != nil {
			err = trx.Rollback()
			if err != nil {
				slog.Error("failed to rollback transaction", "error", err)
			}
			slog.Error("recovered from panic", "error", rec)
			panic(rec)
		}
		if err != nil {
			if rbErr := trx.Rollback(); rbErr != nil {
				slog.Error("failed to rollback transaction", "error", rbErr)
			}
		} else {
			if cErr := trx.Commit(); cErr != nil {
				slog.Error("failed to commit transaction", "error", cErr)
				err = cErr
			}
		}
	}()
	ctx = context.WithValue(ctx, txKey, trx)
	err = fn(ctx)
	return err
}

// ExtractTx - функция хелпер для получения ключа транзакции из контекста
func ExtractTx(ctx context.Context) (*sqlx.Tx, bool) {
	v, ok := ctx.Value(txKey).(*sqlx.Tx)
	if ok {
		return v, true
	}
	return nil, false
}

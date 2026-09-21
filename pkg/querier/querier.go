package querier

import (
	"context"
	"database/sql"
)

// Querier - интерфейс для устранения дублирования в репозиториях.
// Принимает, как *sqlx.DB, так и *sqlx.Tx - репозиторий не знает,
// работает ли он в транзакции или нет
type Querier interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

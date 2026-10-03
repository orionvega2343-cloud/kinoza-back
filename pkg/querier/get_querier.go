package querier

import (
	"context"
	"kinoza-back/pkg/transaction"
)

func GetQuerier(ctx context.Context, base Querier) Querier {
	q := base
	if tx, ok := transaction.ExtractTx(ctx); ok {
		q = tx
	}
	return q
}

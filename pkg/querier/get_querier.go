package repository

import (
	"context"
	"kinoza-back/pkg/querier"
	"kinoza-back/pkg/transaction"
)

func GetQuerier(ctx context.Context, base querier.Querier) querier.Querier {
	q := base
	if tx, ok := transaction.ExtractTx(ctx); ok {
		q = tx
	}
	return q
}

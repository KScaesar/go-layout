package utility

import (
	"context"

	"gorm.io/gorm"
)

// Transaction It abstracts the underlying transaction management,
// allowing the application layer to focus on business logic.
//
// By using context.Context for transaction propagation, this approach
// enhances separation of concerns,
// enabling a clean separation between business logic and data handling.
//
// This abstraction also allows for easy switching of the underlying database in the future,
// without impacting the overall business logic or requiring significant code changes.
type Transaction func(ctx context.Context, flow func(ctxTX context.Context) error) error

func NonTransaction() Transaction {
	return func(ctx context.Context, flow func(context.Context) error) error {
		return flow(ctx)
	}
}

// NewGormTransaction 把 tx *gorm.DB 放在 context.Context 進行參數傳遞,
// 如此一來, 在應用服務層就可以隱藏 tx 物件, 只依賴抽象的 repository,
// 而且在資料層也可以透過 CtxGetGormTransaction 取得 *gorm.DB
func NewGormTransaction(db *gorm.DB) Transaction {
	return func(ctx context.Context, flow func(context.Context) error) error {
		return db.Transaction(func(tx *gorm.DB) error {
			return flow(CtxWithGromTransaction(ctx, db, tx))
		})
	}
}

func CtxWithGromTransaction(ctx context.Context, database *gorm.DB, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, database, tx)
}

func CtxGetGormTransaction(ctx context.Context, database *gorm.DB) *gorm.DB {
	tx, ok := ctx.Value(database).(*gorm.DB)
	if !ok {
		return database
	}
	return tx
}

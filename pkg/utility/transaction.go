package utility

import (
	"context"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CtxWithGormTX(ctx context.Context, database *gorm.DB, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, database, tx)
}

func CtxGetGormTX(ctx context.Context, database *gorm.DB) *gorm.DB {
	tx, ok := ctx.Value(database).(*gorm.DB)
	if !ok {
		return database
	}
	return tx
}

//

type EasyTransaction func(ctx context.Context, flow func(ctxTX context.Context) error) error

func NonEasyTransaction() EasyTransaction {
	return func(ctx context.Context, flow func(context.Context) error) error {
		return flow(ctx)
	}
}

// NewGormEasyTransaction 把 tx *gorm.DB 放在 context.Context 進行參數傳遞,
// 如此一來, 在應用服務層就可以隱藏 tx 物件, 只依賴抽象的 repository,
// 而且在資料層也可以透過 CtxGetGormTX 取得 *gorm.DB
func NewGormEasyTransaction(db *gorm.DB) EasyTransaction {
	return func(ctx context.Context, flow func(context.Context) error) error {
		return db.Transaction(func(tx *gorm.DB) error {
			return flow(CtxWithGormTX(ctx, db, tx))
		})
	}
}

//

// Transaction It abstracts the underlying transaction management,
// allowing the application layer to focus on business logic.
//
// By using context.Context for transaction propagation, this approach
// enhances separation of concerns,
// enabling a clean separation between business logic and data handling.
//
// This abstraction also allows for easy switching of the underlying database in the future,
// without impacting the overall business logic or requiring significant code changes.
type Transaction interface {
	Begin(ctx context.Context) (ctxTX context.Context, err error)
	Commit(ctxTX context.Context) error
	Rollback(ctxTX context.Context) error
}

func NewGormTransaction(db *gorm.DB) Transaction {
	return &gormTX{db: db}
}

type gormTX struct {
	db *gorm.DB
}

func (g *gormTX) Begin(ctx context.Context) (ctxTX context.Context, err error) {
	tx := g.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return CtxWithGormTX(ctx, g.db, tx), nil
}

func (g *gormTX) Commit(ctxTX context.Context) error {
	return CtxGetGormTX(ctxTX, g.db).Commit().Error
}

func (g *gormTX) Rollback(ctxTX context.Context) error {
	return CtxGetGormTX(ctxTX, g.db).Rollback().Error
}

//

// GinGormTransaction
//
// 若 skipPaths 長度為 0，表示所有 path 都會使用 tx
// 若 skipPaths 長度不為 0，表示 skipPaths 中的路徑將不會啟動 ix
func GinGormTransaction(db *gorm.DB, skipPaths []string) gin.HandlerFunc {
	skipMethods := map[string]bool{
		http.MethodHead:    true,
		http.MethodConnect: true,
		http.MethodOptions: true,
		http.MethodTrace:   true,

		http.MethodGet:    false,
		http.MethodPost:   false,
		http.MethodPut:    false,
		http.MethodPatch:  false,
		http.MethodDelete: false,
	}

	return func(c *gin.Context) {
		canSkip := db == nil ||
			skipMethods[c.Request.Method] ||
			(len(skipPaths) != 0 && slices.Contains(skipPaths, c.Request.URL.Path))

		if canSkip {
			c.Next()
			return
		}

		tx := db.Begin()
		err := tx.Error
		if err != nil {
			c.Error(err)
			return
		}

		c.Request = c.Request.WithContext(CtxWithGormTX(c.Request.Context(), db, tx))
		c.Next()

		if len(c.Errors) > 0 {
			err := tx.Rollback().Error
			if err != nil {

			}
			return
		}

		err = tx.Commit().Error
		if err != nil {
			c.Error(err)
			return
		}
	}
}

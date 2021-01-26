package dbprovider

import (
	"context"
	"database/sql"

	"github.com/cenk/backoff"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GormProvider gorm adapter for db provider
type GormProvider interface {
	DBProvider
	GetDB(ctx context.Context) *gorm.DB
	Lock(ctx context.Context, write bool) *gorm.DB
}

// NewGorm new gorm db provider
func NewGorm(c *DBConfig, gormConfig *gorm.Config) (GormProvider, error) {
	var (
		dsn       = c.ConnString()
		db        *gorm.DB
		dialector gorm.Dialector
		err       error
	)

	switch c.Type {
	case Mysql:
		dialector = mysql.Open(dsn)
	case Pg:
		dialector = postgres.Open(dsn)
	case Sqlite:
		dialector = sqlite.Open(dsn)
	}

	bo := backoff.NewExponentialBackOff()
	err = backoff.Retry(func() error {
		db, err = gorm.Open(dialector, gormConfig)
		if err != nil {
			return err
		}

		sqlDB, err := db.DB()
		if err != nil {
			return err
		}

		err = sqlDB.Ping()
		return err
	}, bo)

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	c.ConfigWithDB(sqlDB)

	return &gormDB{db: db}, nil
}

type gormDB struct {
	db *gorm.DB
}

// DB get raw sql database instance
func (db *gormDB) DB() (*sql.DB, error) {
	return db.db.DB()
}

// forUpdate for transactions rollback
func (db *gormDB) Lock(ctx context.Context, write bool) *gorm.DB {
	if write {
		return db.GetDB(ctx).Clauses(clause.Locking{Strength: "UPDATE", Table: clause.Table{Name: clause.CurrentTable}})
	}
	return db.GetDB(ctx).Clauses(clause.Locking{Strength: "SHARE", Table: clause.Table{Name: clause.CurrentTable}})
}

func (db *gormDB) getTx(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(TxKey{}).(*gorm.DB)
	if !ok {
		return tx
	}
	return nil
}

func (db *gormDB) GetDB(ctx context.Context) *gorm.DB {
	if tx := db.getTx(ctx); tx != nil {
		return tx
	}

	return db.db.WithContext(ctx)
}

func (db *gormDB) Begin(ctx context.Context) context.Context {
	tx := db.db.Begin().WithContext(ctx)
	return context.WithValue(ctx, TxKey{}, tx)
}

func (db *gormDB) Commit(ctx context.Context) error {
	tx := db.getTx(ctx)
	if tx == nil {
		return ErrNilTx
	}
	tx.Commit()
	return nil
}

func (db *gormDB) Rollback(ctx context.Context) error {
	tx := db.getTx(ctx)
	if tx == nil {
		return ErrNilTx
	}
	tx.Rollback()
	return nil
}

// ExecuteTx start transaction
func (db *gormDB) ExecuteTx(ctx context.Context, callback func(txCtx context.Context) error) error {
	txCtx := db.Begin(ctx)

	var callbackErr, txErr error

	callbackErr = callback(txCtx)

	if callbackErr != nil {
		txErr = db.Rollback(txCtx)
	} else {
		if txErr = db.Commit(txCtx); txErr != nil {
			txErr = db.Rollback(txCtx)
		}
	}

	if txErr != nil {
		if callbackErr != nil {
			txErr = errors.Wrapf(txErr, "transaction callback error:%+v", callbackErr)
		}
		return txErr
	}
	return callbackErr
}

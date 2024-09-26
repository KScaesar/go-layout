package adapters

import (
	"fmt"

	"gorm.io/driver/mysql"

	"github.com/KScaesar/go-layout/pkg"

	"gorm.io/gorm"
)

func NewMySqlGorm(conf *pkg.MySql) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(conf.DSN()), &gorm.Config{
		Logger:                                   nil,
		NowFunc:                                  nil,
		DryRun:                                   false,
		DisableForeignKeyConstraintWhenMigrating: false,
		IgnoreRelationshipsWhenMigrating:         false,
	})
	if err != nil {
		return nil, fmt.Errorf("connect mysql: %w", err)
	}

	stdDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get stdDB: %w", err)
	}

	err = stdDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping mysql: %w", err)
	}

	if conf.Debug {
		db = db.Debug()
	}

	pkg.Shutdown().AddPriorityShutdownAction(2, "mysql", stdDB.Close)
	return db, nil
}

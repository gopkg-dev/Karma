package gormx

import (
	"context"

	"github.com/google/wire"
	"gorm.io/gorm"
)

// TransSet 注入Trans
var TransSet = wire.NewSet(wire.Struct(new(Trans), "*"))

type Trans struct {
	DB *gorm.DB
}

// TransFunc Define transaction execute function
type TransFunc func(context.Context) error

func (a *Trans) Exec(ctx context.Context, fn TransFunc) error {
	if _, ok := FromTrans(ctx); ok {
		return fn(ctx)
	}
	return a.DB.Transaction(func(db *gorm.DB) error {
		return fn(NewTrans(ctx, db))
	})
}

func ExecTrans(ctx context.Context, db *gorm.DB, fn TransFunc) error {
	transModel := &Trans{DB: db}
	return transModel.Exec(ctx, fn)
}

func ExecTransWithLock(ctx context.Context, db *gorm.DB, fn TransFunc) error {
	if !FromRowLock(ctx) {
		ctx = NewRowLock(ctx)
	}
	return ExecTrans(ctx, db, fn)
}

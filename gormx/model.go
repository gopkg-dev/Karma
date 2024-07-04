package gormx

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Model base model
type Model struct {
	ID        string         `gorm:"column:id;size:20;primaryKey"` // Unique ID
	CreatedAt time.Time      `gorm:"column:created_at;index;"`     // Create time
	UpdatedAt time.Time      `gorm:"column:updated_at;index;"`     // Update time
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`      // Delete time
}

// GetDB Get gorm.DB from context
func GetDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	db := defDB
	if tdb, ok := FromTrans(ctx); ok {
		db = tdb
	}
	if FromRowLock(ctx) {
		db = db.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	return db.WithContext(ctx)
}

// GetDBWithModel Get gorm.DB.Model from context
func GetDBWithModel(ctx context.Context, defDB *gorm.DB, m interface{}) *gorm.DB {
	return GetDB(ctx, defDB).Model(m)
}

package logging

import (
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type Logger struct {
	ID        string    `gorm:"size:20;primaryKey;" json:"id"`  // Unique ID
	Level     string    `gorm:"size:20;index;" json:"level"`    // Log level
	TraceID   string    `gorm:"size:64;index;" json:"trace_id"` // Trace ID
	Tag       string    `gorm:"size:50;index;" json:"tag"`      // Log tag
	Message   string    `gorm:"size:1024;" json:"message"`      // Log message
	Stack     string    `gorm:"type:text;" json:"stack"`        // Error stack
	Data      string    `gorm:"type:text;" json:"data"`         // Log data
	UserID    string    `gorm:"size:20;index;" json:"user_id"`  // User ID
	CreatedAt time.Time `gorm:"index;" json:"created_at"`       // Create time
}

func NewGormHook(db *gorm.DB) *GormHook {
	err := db.AutoMigrate(new(Logger))
	if err != nil {
		panic(err)
	}

	return &GormHook{
		db: db,
	}
}

// GormHook Gorm Logger Hook
type GormHook struct {
	db *gorm.DB
}

func (h *GormHook) Exec(extra map[string]string, data []byte) error {
	msg := &Logger{
		ID: xid.New().String(),
	}

	logData := make(map[string]interface{})
	if err := jsoniter.Unmarshal(data, &logData); err != nil {
		return err
	}

	for key, value := range logData {
		switch key {
		case "ts":
			msg.CreatedAt = time.UnixMilli(int64(value.(float64)))
		case "msg":
			msg.Message = value.(string)
		case "tag":
			msg.Tag = value.(string)
		case "trace_id":
			msg.TraceID = value.(string)
		case "user_id":
			msg.UserID = value.(string)
		case "level":
			msg.Level = value.(string)
		case "stack":
			msg.Stack = value.(string)
		}
	}

	delete(logData, "caller")

	for k, v := range extra {
		logData[k] = v
	}

	if len(data) > 0 {
		buf, _ := jsoniter.Marshal(data)
		msg.Data = string(buf)
	}

	return h.db.Create(msg).Error
}

func (h *GormHook) Close() error {
	db, err := h.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

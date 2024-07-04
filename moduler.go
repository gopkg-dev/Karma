package karma

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// Moduler 是一个接口，定义了模块的方法。
type Moduler interface {
	// Stringer 表示模块的名称
	fmt.Stringer

	// Init 执行初始化工作。
	Init(context.Context) error

	// AutoMigrate 执行数据库迁移。
	AutoMigrate(context.Context) error

	// RegisterRoutes 向 Fiber 路由器添加路由。
	RegisterRoutes(context.Context, fiber.Router)

	// Release 释放模块资源。
	Release(context.Context) error
}

// Module 是一个空结构体，实现了 Moduler 接口
// 可以作为 Moduler 嵌入到自定义结构体中。
type Module struct{}

// String 返回模块的名称。
func (Module) String() string { return "anonymous" }

// AutoMigrate 执行数据库迁移。
func (Module) AutoMigrate(context.Context) error { return nil }

// Init 执行初始化工作。
func (Module) Init(context.Context) error { return nil }

// RegisterRoutes 向 Fiber 路由器添加路由。
func (Module) RegisterRoutes(ctx context.Context, router fiber.Router) {}

// Release 释放模块资源。
func (Module) Release(ctx context.Context) error { return nil }

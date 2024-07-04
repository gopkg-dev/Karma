package karma

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockModule struct {
	Module
}

// go test -run Test_Moduler_Embed_Empty_Module -race
func Test_Moduler_Embed_Empty_Module(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	module := mockModule{Module{}}

	assert.Implements(t, (*Moduler)(nil), module)

	assert.Equal(t, "anonymous", module.String())

	assert.Nil(t, module.Init(ctx))
	assert.Nil(t, module.AutoMigrate(ctx))

	module.RegisterRoutes(ctx, nil)

	assert.Nil(t, module.Release(ctx))
}

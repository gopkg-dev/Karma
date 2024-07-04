package jwtx_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gopkg-dev/karma/jwtx"
)

func TestAuth(t *testing.T) {
	cache := jwtx.NewMemoryCache(jwtx.MemoryConfig{CleanupInterval: time.Second})

	store := jwtx.NewStoreWithCache(cache)
	ctx := context.Background()
	jwtAuth := jwtx.New(store)

	userID := "test"
	token, err := jwtAuth.GenerateToken(ctx, userID)
	assert.Nil(t, err)
	assert.NotNil(t, token)

	id, err := jwtAuth.ParseSubject(ctx, token.GetAccessToken())
	assert.Nil(t, err)
	assert.Equal(t, userID, id)

	err = jwtAuth.DestroyToken(ctx, token.GetAccessToken())
	assert.Nil(t, err)

	id, err = jwtAuth.ParseSubject(ctx, token.GetAccessToken())
	assert.NotNil(t, err)
	assert.EqualError(t, err, jwtx.ErrTokenInvalid.Error())
	assert.Empty(t, id)

	err = jwtAuth.Release(ctx)
	assert.Nil(t, err)
}

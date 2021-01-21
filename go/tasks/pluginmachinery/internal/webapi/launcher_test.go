package webapi

import (
	"context"
	"fmt"
	"testing"

	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/core/mocks"
	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/webapi"
	mocks2 "github.com/lyft/flytestdlib/cache/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_launch(t *testing.T) {
	t.Run("Successful launch", func(t *testing.T) {
		ctx := context.Background()
		tCtx := &mocks.TaskExecutionContext{}
		meta := &mocks.TaskExecutionMetadata{}
		taskID := &mocks.TaskExecutionID{}
		taskID.OnGetGeneratedName().Return("my-id")
		meta.OnGetTaskExecutionID().Return(taskID)
		tCtx.OnTaskExecutionMetadata().Return(meta)

		c := &mocks2.AutoRefresh{}
		s := State{}
		c.OnGetOrCreate("my-id", CacheItem{State: s}).Return(CacheItem{State: s}, nil)

		plgn := newPluginWithProperties(webapi.PluginConfig{})
		plgn.OnCreate(ctx, tCtx).Return("abc", nil)
		newS, phaseInfo, err := launch(ctx, plgn, tCtx, c, &s)
		assert.NoError(t, err)
		assert.NotNil(t, newS)
		assert.NotNil(t, phaseInfo)
	})

	t.Run("Failed to create resource", func(t *testing.T) {
		ctx := context.Background()
		tCtx := &mocks.TaskExecutionContext{}
		meta := &mocks.TaskExecutionMetadata{}
		taskID := &mocks.TaskExecutionID{}
		taskID.OnGetGeneratedName().Return("my-id")
		meta.OnGetTaskExecutionID().Return(taskID)
		tCtx.OnTaskExecutionMetadata().Return(meta)

		c := &mocks2.AutoRefresh{}
		s := State{}
		c.OnGetOrCreate("my-id", CacheItem{State: s}).Return(CacheItem{State: s}, nil)

		plgn := newPluginWithProperties(webapi.PluginConfig{})
		plgn.OnCreate(ctx, tCtx).Return("", fmt.Errorf("error creating"))
		_, _, err := launch(ctx, plgn, tCtx, c, &s)
		assert.Error(t, err)
	})

	t.Run("Failed to cache", func(t *testing.T) {
		ctx := context.Background()
		tCtx := &mocks.TaskExecutionContext{}
		meta := &mocks.TaskExecutionMetadata{}
		taskID := &mocks.TaskExecutionID{}
		taskID.OnGetGeneratedName().Return("my-id")
		meta.OnGetTaskExecutionID().Return(taskID)
		tCtx.OnTaskExecutionMetadata().Return(meta)

		c := &mocks2.AutoRefresh{}
		s := State{}
		c.OnGetOrCreate("my-id", CacheItem{State: s}).Return(CacheItem{State: s}, fmt.Errorf("failed to cache"))

		plgn := newPluginWithProperties(webapi.PluginConfig{})
		plgn.OnCreate(ctx, tCtx).Return("my-id", nil)
		_, _, err := launch(ctx, plgn, tCtx, c, &s)
		assert.Error(t, err)
	})
}

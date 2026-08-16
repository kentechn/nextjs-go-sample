package todo_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kentechn/nextjs-go-sample/apps/api/internal/todo"
)

func TestStoreCreateAndGet(t *testing.T) {
	t.Parallel()

	store := todo.NewStore()
	created := store.Create("write the spec")

	got, err := store.Get(created.ID)
	require.NoError(t, err)
	assert.Equal(t, "write the spec", got.Title)
	assert.False(t, got.Done)
}

func TestStoreListFiltersByDone(t *testing.T) {
	t.Parallel()

	store := todo.NewStore()
	store.Create("first")
	store.Create("second")

	done := false
	assert.Len(t, store.List(&done), 2)

	done = true
	assert.Empty(t, store.List(&done))
	assert.Len(t, store.List(nil), 2)
}

func TestStoreDelete(t *testing.T) {
	t.Parallel()

	store := todo.NewStore()
	created := store.Create("temporary")

	require.NoError(t, store.Delete(created.ID))
	assert.ErrorIs(t, store.Delete(created.ID), todo.ErrNotFound)
	assert.ErrorIs(t, store.Delete(uuid.New()), todo.ErrNotFound)
}

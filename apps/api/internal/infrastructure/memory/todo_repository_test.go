package memory_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domain "github.com/kentechn/nextjs-go-sample/apps/api/internal/domain/todo"
	"github.com/kentechn/nextjs-go-sample/apps/api/internal/infrastructure/memory"
)

var baseTime = time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)

func item(t *testing.T, title string, offset time.Duration) domain.Todo {
	t.Helper()

	created, err := domain.New(uuid.New(), title, baseTime.Add(offset))
	require.NoError(t, err)

	return created
}

func TestCreateAndFindByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := memory.NewTodoRepository()
	created := item(t, "write the spec", 0)
	require.NoError(t, repo.Create(ctx, created))

	got, err := repo.FindByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created, got)
}

func TestCreateRejectsDuplicateID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := memory.NewTodoRepository()
	created := item(t, "write the spec", 0)

	require.NoError(t, repo.Create(ctx, created))
	assert.ErrorIs(t, repo.Create(ctx, created), domain.ErrAlreadyExists)
}

func TestListIsSortedAndFiltered(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := memory.NewTodoRepository()
	second := item(t, "second", time.Minute)
	first := item(t, "first", 0)
	require.NoError(t, repo.Create(ctx, second))
	require.NoError(t, repo.Create(ctx, first))

	all, err := repo.List(ctx, domain.Filter{})
	require.NoError(t, err)
	require.Len(t, all, 2)
	assert.Equal(t, []string{"first", "second"}, []string{all[0].Title, all[1].Title})

	done := true
	completed, err := repo.List(ctx, domain.Filter{Done: &done})
	require.NoError(t, err)
	assert.Empty(t, completed)
}

func TestDelete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := memory.NewTodoRepository()
	created := item(t, "temporary", 0)
	require.NoError(t, repo.Create(ctx, created))

	require.NoError(t, repo.Delete(ctx, created.ID))
	assert.ErrorIs(t, repo.Delete(ctx, created.ID), domain.ErrNotFound)
	_, err := repo.FindByID(ctx, uuid.New())
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

package todo_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domain "github.com/kentechn/nextjs-go-sample/apps/api/internal/domain/todo"
	usecase "github.com/kentechn/nextjs-go-sample/apps/api/internal/usecase/todo"
)

var (
	fixedID   = uuid.MustParse("2f1c3f8a-0b1d-4d5e-8a7b-9c0d1e2f3a4b")
	fixedTime = time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)
	errBroken = errors.New("repository is broken")
)

// fakeRepository is a stand-in for the infrastructure layer: the use cases only
// know the domain interface, so no storage is needed to test them.
type fakeRepository struct {
	items   []domain.Todo
	failErr error
}

func (r *fakeRepository) List(_ context.Context, filter domain.Filter) ([]domain.Todo, error) {
	if r.failErr != nil {
		return nil, r.failErr
	}

	out := make([]domain.Todo, 0, len(r.items))
	for _, item := range r.items {
		if filter.Matches(item) {
			out = append(out, item)
		}
	}

	return out, nil
}

func (r *fakeRepository) FindByID(_ context.Context, id uuid.UUID) (domain.Todo, error) {
	if r.failErr != nil {
		return domain.Todo{}, r.failErr
	}
	for _, item := range r.items {
		if item.ID == id {
			return item, nil
		}
	}

	return domain.Todo{}, domain.ErrNotFound
}

func (r *fakeRepository) Create(_ context.Context, item domain.Todo) error {
	if r.failErr != nil {
		return r.failErr
	}
	r.items = append(r.items, item)

	return nil
}

func (r *fakeRepository) Delete(_ context.Context, id uuid.UUID) error {
	if r.failErr != nil {
		return r.failErr
	}
	for i, item := range r.items {
		if item.ID == id {
			r.items = append(r.items[:i], r.items[i+1:]...)

			return nil
		}
	}

	return domain.ErrNotFound
}

func newUseCase(repo domain.Repository) *usecase.UseCase {
	return usecase.New(
		repo,
		usecase.WithClock(func() time.Time { return fixedTime }),
		usecase.WithIDGenerator(func() uuid.UUID { return fixedID }),
	)
}

func TestCreateUsesInjectedClockAndID(t *testing.T) {
	t.Parallel()

	repo := &fakeRepository{}
	item, err := newUseCase(repo).Create(context.Background(), usecase.CreateInput{Title: "ship it"})

	require.NoError(t, err)
	assert.Equal(t, fixedID, item.ID)
	assert.Equal(t, fixedTime, item.CreatedAt)
	assert.Len(t, repo.items, 1)
}

// Domain rules are enforced before anything is stored.
func TestCreateRejectsInvalidTitle(t *testing.T) {
	t.Parallel()

	repo := &fakeRepository{}
	useCase := newUseCase(repo)

	_, err := useCase.Create(context.Background(), usecase.CreateInput{Title: "   "})
	assert.ErrorIs(t, err, domain.ErrEmptyTitle)

	_, err = useCase.Create(context.Background(), usecase.CreateInput{
		Title: strings.Repeat("a", domain.TitleMaxLength+1),
	})
	assert.ErrorIs(t, err, domain.ErrTitleTooLong)
	assert.Empty(t, repo.items)
}

func TestListFiltersByDone(t *testing.T) {
	t.Parallel()

	repo := &fakeRepository{items: []domain.Todo{
		{ID: uuid.New(), Title: "open", Done: false},
		{ID: uuid.New(), Title: "done", Done: true},
	}}
	useCase := newUseCase(repo)
	done := true

	all, err := useCase.List(context.Background(), usecase.ListInput{})
	require.NoError(t, err)
	assert.Len(t, all, 2)

	completed, err := useCase.List(context.Background(), usecase.ListInput{Done: &done})
	require.NoError(t, err)
	require.Len(t, completed, 1)
	assert.Equal(t, "done", completed[0].Title)
}

func TestGetAndDeleteMissingTodo(t *testing.T) {
	t.Parallel()

	useCase := newUseCase(&fakeRepository{})

	_, err := useCase.Get(context.Background(), uuid.New())
	assert.ErrorIs(t, err, domain.ErrNotFound)
	assert.ErrorIs(t, useCase.Delete(context.Background(), uuid.New()), domain.ErrNotFound)
}

// Repository failures are wrapped, not swallowed.
func TestRepositoryErrorIsPropagated(t *testing.T) {
	t.Parallel()

	useCase := newUseCase(&fakeRepository{failErr: errBroken})

	_, err := useCase.List(context.Background(), usecase.ListInput{})
	assert.ErrorIs(t, err, errBroken)
}

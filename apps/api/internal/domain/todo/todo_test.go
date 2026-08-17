package todo_test

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kentechn/nextjs-go-sample/apps/api/internal/domain/todo"
)

var (
	testID   = uuid.MustParse("2f1c3f8a-0b1d-4d5e-8a7b-9c0d1e2f3a4b")
	testTime = time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
)

func TestNewTrimsTitle(t *testing.T) {
	t.Parallel()

	item, err := todo.New(testID, "  write the spec  ", testTime)
	require.NoError(t, err)
	assert.Equal(t, "write the spec", item.Title)
	assert.Equal(t, testID, item.ID)
	assert.Equal(t, testTime, item.CreatedAt)
	assert.False(t, item.Done)
}

func TestNewRejectsInvalidTitle(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		title string
		want  error
	}{
		"empty":           {title: "", want: todo.ErrEmptyTitle},
		"whitespace only": {title: " \t\n", want: todo.ErrEmptyTitle},
		"one rune too long": {
			title: strings.Repeat("あ", todo.TitleMaxLength+1),
			want:  todo.ErrTitleTooLong,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := todo.New(testID, test.title, testTime)
			assert.ErrorIs(t, err, test.want)
		})
	}
}

// The maximum length is counted in runes, not bytes, so a title of 200
// multi-byte characters is still valid.
func TestNewAcceptsMaxLengthTitle(t *testing.T) {
	t.Parallel()

	_, err := todo.New(testID, strings.Repeat("あ", todo.TitleMaxLength), testTime)
	require.NoError(t, err)
}

func TestFilterMatches(t *testing.T) {
	t.Parallel()

	done := todo.Todo{Done: true}
	open := todo.Todo{Done: false}
	value := true

	assert.True(t, todo.Filter{}.Matches(done))
	assert.True(t, todo.Filter{}.Matches(open))
	assert.True(t, todo.Filter{Done: &value}.Matches(done))
	assert.False(t, todo.Filter{Done: &value}.Matches(open))
}

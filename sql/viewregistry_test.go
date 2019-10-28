package sql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	dbName   = "db"
	viewName = "myview"
	mockView = NewView(viewName, nil)
)

// Tests the creation of an empty ViewRegistry with no views registered.
func TestNewViewRegistry(t *testing.T) {
	require := require.New(t)

	registry := NewViewRegistry()
	require.Equal(0, len(registry.AllViews()))
}

// Tests that registering a non-existing view succeeds.
func TestRegisterNonExistingView(t *testing.T) {
	require := require.New(t)

	registry := NewViewRegistry()

	err := registry.Register(dbName, mockView)
	require.NoError(err)
	require.Equal(1, len(registry.AllViews()))

	actualView, err := registry.View(dbName, viewName)
	require.NoError(err)
	require.Equal(mockView, *actualView)
}

// Tests that registering an existing view fails.
func TestRegisterExistingVIew(t *testing.T) {
	require := require.New(t)

	registry := NewViewRegistry()

	err := registry.Register(dbName, mockView)
	require.NoError(err)
	require.Equal(1, len(registry.AllViews()))

	err = registry.Register(dbName, mockView)
	require.Error(err)
	require.True(ErrExistingView.Is(err))
}

// Tests that deleting an existing view succeeds.
func TestDeleteExistingView(t *testing.T) {
	require := require.New(t)

	registry := NewViewRegistry()

	err := registry.Register(dbName, mockView)
	require.NoError(err)
	require.Equal(1, len(registry.AllViews()))

	err = registry.Delete(dbName, viewName)
	require.NoError(err)
	require.Equal(0, len(registry.AllViews()))
}

// Tests that deleting a non-existing view fails.
func TestDeleteNonExistingView(t *testing.T) {
	require := require.New(t)

	registry := NewViewRegistry()

	err := registry.Delete("random", "randomer")
	require.Error(err)
	require.True(ErrNonExistingView.Is(err))
}

// Tests that retrieving an existing view succeeds and that the view returned
// is the correct one.
func TestGetExistingView(t *testing.T) {
	require := require.New(t)

	registry := NewViewRegistry()

	err := registry.Register(dbName, mockView)
	require.NoError(err)
	require.Equal(1, len(registry.AllViews()))

	actualView, err := registry.View(dbName, viewName)
	require.NoError(err)
	require.Equal(mockView, *actualView)
}

// Tests that retrieving a non-existing view fails.
func TestGetNonExistingView(t *testing.T) {
	require := require.New(t)

	registry := NewViewRegistry()

	actualView, err := registry.View(dbName, viewName)
	require.Error(err)
	require.Nil(actualView)
	require.True(ErrNonExistingView.Is(err))
}

// Tests that retrieving the views registered under a database succeeds,
// returning the list of all the correct views.
func TestViewsInDatabase(t *testing.T) {
	require := require.New(t)

	registry := NewViewRegistry()

	databases := []struct {
		name     string
		numViews int
	}{
		{"db0", 0},
		{"db1", 5},
		{"db2", 10},
	}

	for _, db := range databases {
		for i := 0; i < db.numViews; i++ {
			view := NewView(viewName+string(i), nil)
			err := registry.Register(db.name, view)
			require.NoError(err)
		}

		views := registry.ViewsInDatabase(db.name)
		require.Equal(db.numViews, len(views))
	}
}

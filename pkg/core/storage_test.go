package core

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorageTerm(t *testing.T) {
	path := "./storage/0"
	term0 := Term{
		Id: "1",
		No: 0,
	}
	first0 := Value("1")
	term1 := Term{
		Id: "1",
		No: 1,
	}

	storage, err := NewStorage(path)
	require.NoError(t, err)

	err = storage.Create(term0, first0)
	require.NoError(t, err)

	err = storage.PutTerm(term1)
	require.NoError(t, err)

	term, err := storage.GetTerm()
	require.NoError(t, err)
	require.Equal(t, term1, term)

	err = storage.Delete()
	require.NoError(t, err)

	storage.Stop()

	err = os.Remove(path)
	require.NoError(t, err)
}

func TestStorageFirst(t *testing.T) {
	path := "./storage/0"
	term0 := Term{
		Id: "1",
		No: 0,
	}
	first0 := Value("1")
	first1 := Value("2")

	storage, err := NewStorage(path)
	require.NoError(t, err)

	err = storage.Create(term0, first0)
	require.NoError(t, err)

	err = storage.PutFirst(first1)
	require.NoError(t, err)

	first, err := storage.GetFirst()
	require.NoError(t, err)
	require.Equal(t, first1, first)

	err = storage.Delete()
	require.NoError(t, err)

	storage.Stop()

	err = os.Remove(path)
	require.NoError(t, err)
}

func TestStorageValue(t *testing.T) {
	path := "./storage/0"
	term0 := Term{
		Id: "1",
		No: 0,
	}
	first0 := Value("1")
	value0 := Value("1")

	storage, err := NewStorage(path)
	require.NoError(t, err)

	err = storage.Create(term0, first0)
	require.NoError(t, err)

	err = storage.PutValue(value0)
	require.NoError(t, err)

	value, err := storage.GetValue()
	require.NoError(t, err)
	require.Equal(t, value0, value)

	err = storage.Delete()
	require.NoError(t, err)

	storage.Stop()

	err = os.Remove(path)
	require.NoError(t, err)
}

func TestStorageVote(t *testing.T) {
	path := "./storage/0"
	term0 := Term{
		Id: "1",
		No: 0,
	}
	first0 := Value("1")
	vote0 := Vote{
		Term:  term0,
		Value: "1",
	}

	storage, err := NewStorage(path)
	require.NoError(t, err)

	err = storage.Create(term0, first0)
	require.NoError(t, err)

	err = storage.PutVote(vote0)
	require.NoError(t, err)

	vote, err := storage.GetVote()
	require.NoError(t, err)
	require.Equal(t, vote0, vote)

	err = storage.Delete()
	require.NoError(t, err)

	storage.Stop()

	err = os.Remove(path)
	require.NoError(t, err)
}

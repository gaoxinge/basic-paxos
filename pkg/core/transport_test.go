package core

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gaoxinge/basic-paxos/pkg/util/log"
)

func TestTransport1(t *testing.T) {
	ids := []string{"1", "2", "3"}
	paths := []string{"./storage/1", "./storage/2", "./storage/3"}
	firsts := []string{"1", "2", "3"}
	prepareNum := 2
	voteNum := 2

	testServer, err := NewTestServer(ids, paths, firsts, prepareNum, voteNum)
	require.NoError(t, err)

	for _, id := range testServer.Ids {
		p := testServer.Paxos[id]
		log.Info().Msgf("[HANDLE] paxos %s handle timeout", id)
		err = p.Handle(Message{Type: Timeout})
		require.NoError(t, err)
		require.True(t, testServer.Check())
	}

	for i := 0; i < 5000; i++ {
		err = testServer.Handle()
		require.NoError(t, err)
		require.True(t, testServer.Check())
	}
}

func TestTransport2(t *testing.T) {
	ids := []string{"1", "2", "3", "4", "5"}
	paths := []string{"./storage/1", "./storage/2", "./storage/3", "./storage/4", "./storage/5"}
	firsts := []string{"1", "2", "3", "4", "5"}
	prepareNum := 3
	voteNum := 3

	testServer, err := NewTestServer(ids, paths, firsts, prepareNum, voteNum)
	require.NoError(t, err)

	for _, id := range testServer.Ids {
		p := testServer.Paxos[id]
		log.Info().Msgf("[HANDLE] paxos %s handle timeout", id)
		err = p.Handle(Message{Type: Timeout})
		require.NoError(t, err)
		require.True(t, testServer.Check())
	}

	for i := 0; i < 20000; i++ {
		err = testServer.Handle()
		require.NoError(t, err)
		require.True(t, testServer.Check())
	}
}

func TestTransport3(t *testing.T) {
	ids := []string{"1", "2", "3", "4", "5"}
	paths := []string{"./storage/1", "./storage/2", "./storage/3", "./storage/4", "./storage/5"}
	firsts := []string{"1", "2", "3", "4", "5"}
	prepareNum := 2
	voteNum := 4

	testServer, err := NewTestServer(ids, paths, firsts, prepareNum, voteNum)
	require.NoError(t, err)

	for _, id := range testServer.Ids {
		p := testServer.Paxos[id]
		log.Info().Msgf("[HANDLE] paxos %s handle timeout", id)
		err = p.Handle(Message{Type: Timeout})
		require.NoError(t, err)
		require.True(t, testServer.Check())
	}

	for i := 0; i < 20000; i++ {
		err = testServer.Handle()
		require.NoError(t, err)
		require.True(t, testServer.Check())
	}
}

func TestTransport4(t *testing.T) {
	ids := []string{"1", "2", "3", "4", "5"}
	paths := []string{"./storage/1", "./storage/2", "./storage/3", "./storage/4", "./storage/5"}
	firsts := []string{"1", "2", "3", "4", "5"}
	prepareNum := 4
	voteNum := 2

	testServer, err := NewTestServer(ids, paths, firsts, prepareNum, voteNum)
	require.NoError(t, err)

	for _, id := range testServer.Ids {
		p := testServer.Paxos[id]
		log.Info().Msgf("[HANDLE] paxos %s handle timeout", id)
		err = p.Handle(Message{Type: Timeout})
		require.NoError(t, err)
		require.True(t, testServer.Check())
	}

	for i := 0; i < 20000; i++ {
		err = testServer.Handle()
		require.NoError(t, err)
		require.True(t, testServer.Check())
	}
}

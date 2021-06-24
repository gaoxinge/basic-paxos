package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTerm(t *testing.T) {
	term1 := Term{
		Id: "1",
		No: 1,
	}
	term2 := Term{
		Id: "1",
		No: 2,
	}
	require.Equal(t, -1, term1.Compare(term2))
	require.Equal(t, 1, term2.Compare(term1))

	term3 := Term{
		Id: "1",
		No: 1,
	}
	term4 := Term{
		Id: "2",
		No: 2,
	}
	require.Equal(t, -1, term3.Compare(term4))
	require.Equal(t, 1, term4.Compare(term3))

	term5 := Term{
		Id: "1",
		No: 1,
	}
	require.Equal(t, 0, term5.Compare(term5))

	term := Term{
		Id: "1",
		No: 1,
	}
	require.Equal(t, "Term{Id: 1, No: 1}", term.String())
}

func TestValue(t *testing.T) {
	value0 := Value("1")
	value1 := Value("")
	require.Equal(t, false, value0.IsNull())
	require.Equal(t, true, value1.IsNull())

	value := Value("1")
	require.Equal(t, "1", value.String())
}

func TestVote(t *testing.T) {
	term := Term{
		Id: "1",
		No: 1,
	}
	vote0 := Vote{
		Term:  term,
		Value: "1",
	}
	vote1 := Vote{
		Term:  term,
		Value: "",
	}
	require.Equal(t, false, vote0.IsNull())
	require.Equal(t, true, vote1.IsNull())

	vote := Vote{
		Term:  term,
		Value: "1",
	}
	require.Equal(t, "Vote{Term: Term{Id: 1, No: 1}, Value: 1}", vote.String())
}

func TestMessage(t *testing.T) {
	term := Term{
		Id: "1",
		No: 1,
	}
	message := Message{
		Term: term,
		Type: Timeout,
		From: "1",
		To:   "1",
		Ok:   true,
	}
	require.Equal(t, "Message{Term: Term{Id: 1, No: 1}, Type: Timeout, From: 1, To: 1, Ok: true, Content: map[]}", message.String())
}

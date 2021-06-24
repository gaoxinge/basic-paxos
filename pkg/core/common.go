package core

import (
	"fmt"
)

type Term struct {
	Id string `json:"id"`
	No int    `json:"no"`
}

func (term Term) Compare(other Term) int {
	if term.No < other.No {
		return -1
	}

	if term.No > other.No {
		return 1
	}

	if term.Id < other.Id {
		return -1
	}

	if term.Id > other.Id {
		return 1
	}

	return 0
}

func (term Term) String() string {
	return fmt.Sprintf("Term{Id: %s, No: %d}", term.Id, term.No)
}

type Value string

func (value Value) IsNull() bool {
	return len(value) == 0
}

func (value Value) String() string {
	return string(value)
}

type Vote struct {
	Term  Term  `json:"term"`
	Value Value `json:"value"`
}

func (vote Vote) String() string {
	return fmt.Sprintf("Vote{Term: %s, Value: %s}", vote.Term, vote.Value)
}

func (vote Vote) IsNull() bool {
	return vote.Value.IsNull()
}

const (
	Timeout = iota
	PrepareRequest
	PrepareResponse
	VoteRequest
	VoteResponse
	Learn
)

type Message struct {
	Term    Term                   `json:"term"`
	Type    int                    `json:"type"`
	From    string                 `json:"from"`
	To      string                 `json:"to"`
	Ok      bool                   `json:"ok"`
	Content map[string]interface{} `json:"content"`
}

func (message Message) String() string {
	t := ""
	switch message.Type {
	case Timeout:
		t = "Timeout"
	case PrepareRequest:
		t = "PrepareRequest"
	case PrepareResponse:
		t = "PrepareResponse"
	case VoteRequest:
		t = "VoteRequest"
	case VoteResponse:
		t = "VoteResponse"
	case Learn:
		t = "Learn"
	}
	return fmt.Sprintf("Message{Term: %s, Type: %s, From: %s, To: %s, Ok: %t, Content: %s}", message.Term, t, message.From, message.To, message.Ok, message.Content)
}

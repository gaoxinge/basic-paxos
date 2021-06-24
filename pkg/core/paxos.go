package core

import (
	"github.com/gaoxinge/basic-paxos/pkg/util/log"
)

const (
	NullRound = iota
	PrepareRound
	VoteRound
)

type Client interface {
	Handle(message Message) error
}

type Server interface {
	Run()
	Stop()
}

type Paxos struct {
	Self   string
	Others map[string]Client

	Path    string
	Storage *Storage

	PrepareNum int
	VoteNum    int

	Term      Term
	TermRound int
	TermAcks  int
	TermVote  Vote

	First Value
	Value Value
	Vote  Vote
}

func NewPaxos(self string, others map[string]Client, path string, prepareNum int, voteNum int, first Value) (*Paxos, error) {
	var err error

	storage, err := NewStorage(path)
	if err != nil {
		return nil, err
	}
	err = storage.Create(Term{
		Id: self,
		No: 1,
	}, first)
	if err != nil {
		return nil, err
	}

	p := Paxos{
		Self:   self,
		Others: others,

		Path:    path,
		Storage: storage,

		PrepareNum: prepareNum,
		VoteNum:    voteNum,
	}

	err = p.getTerm()
	if err != nil {
		return nil, err
	}

	err = p.getFirst()
	if err != nil {
		return nil, err
	}

	err = p.getValue()
	if err != nil {
		return nil, err
	}

	err = p.getVote()
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (p *Paxos) putTerm(term Term) error {
	err := p.Storage.PutTerm(term)
	if err != nil {
		return err
	}
	p.Term = term
	return nil
}

func (p *Paxos) getTerm() error {
	term, err := p.Storage.GetTerm()
	if err != nil {
		return err
	}
	p.Term = term
	return nil
}

func (p *Paxos) putFirst(first Value) error {
	err := p.Storage.PutFirst(first)
	if err != nil {
		return err
	}
	p.First = first
	return nil
}

func (p *Paxos) getFirst() error {
	first, err := p.Storage.GetFirst()
	if err != nil {
		return err
	}
	p.First = first
	return nil
}

func (p *Paxos) putValue(value Value) error {
	err := p.Storage.PutValue(value)
	if err != nil {
		return err
	}
	p.Value = value
	return nil
}

func (p *Paxos) getValue() error {
	value, err := p.Storage.GetValue()
	if err != nil {
		return err
	}
	p.Value = value
	return nil
}

func (p *Paxos) putVote(vote Vote) error {
	err := p.Storage.PutVote(vote)
	if err != nil {
		return err
	}
	p.Vote = vote
	return nil
}

func (p *Paxos) getVote() error {
	vote, err := p.Storage.GetVote()
	if err != nil {
		return err
	}
	p.Vote = vote
	return nil
}

func (p *Paxos) Handle(m Message) error {
	switch m.Type {
	case Timeout:
		return p.handleTimeout(m)
	case PrepareRequest:
		return p.handlePrepareRequest(m)
	case PrepareResponse:
		return p.handlePrepareResponse(m)
	case VoteRequest:
		return p.handleVoteRequest(m)
	case VoteResponse:
		return p.handleVoteResponse(m)
	case Learn:
		return p.handleLearn(m)
	}
	return nil
}

func (p *Paxos) handleTimeout(m Message) error {
	var err error

	if !p.Value.IsNull() {
		return nil
	}

	err = p.putTerm(Term{
		Id: p.Self,
		No: p.Term.No + 1,
	})
	if err != nil {
		return err
	}
	p.TermRound = PrepareRound
	p.TermAcks = 1
	p.TermVote = p.Vote

	for other, client := range p.Others {
		m = Message{
			Term: p.Term,
			Type: PrepareRequest,
			From: p.Self,
			To:   other,
		}
		err = client.Handle(m)
		if err != nil {
			log.Warn().Err(err).Msg("handle timeout with error")
		}
	}

	return nil
}

func (p *Paxos) handlePrepareRequest(m Message) error {
	var err error

	compare := m.Term.Compare(p.Term)
	if compare == -1 {
		m = Message{
			Term: p.Term,
			Type: PrepareResponse,
			From: m.To,
			To:   m.From,
			Ok:   false,
		}
		goto end
	}

	if compare == 1 {
		err = p.putTerm(m.Term)
		if err != nil {
			return err
		}
	}

	m = Message{
		Term:    p.Term,
		Type:    PrepareResponse,
		From:    m.To,
		To:      m.From,
		Ok:      true,
		Content: map[string]interface{}{"vote": p.Vote},
	}

end:
	client := p.Others[m.To]
	err = client.Handle(m)
	if err != nil {
		log.Warn().Err(err).Msg("handle prepare request with error")
	}
	return nil
}

func (p *Paxos) handlePrepareResponse(m Message) error {
	var err error

	compare := m.Term.Compare(p.Term)
	if compare == -1 {
		return nil
	}

	if compare == 1 {
		return p.putTerm(m.Term)
	}

	if compare == 0 && p.TermRound == PrepareRound && m.Ok {
		p.TermAcks++
		vote := m.Content["vote"].(Vote)
		if !vote.IsNull() && (p.TermVote.IsNull() || p.TermVote.Term.Compare(vote.Term) == -1) {
			p.TermVote = vote
		}
		if p.TermAcks == p.PrepareNum {
			p.TermRound = VoteRound
			p.TermAcks = 1

			var value Value
			if !p.TermVote.IsNull() {
				value = p.TermVote.Value
			} else if !p.First.IsNull() {
				value = p.First
			} else {
				return nil
			}

			err = p.putVote(Vote{
				Term:  p.Term,
				Value: value,
			})
			if err != nil {
				return err
			}

			for other, client := range p.Others {
				m = Message{
					Term:    p.Term,
					Type:    VoteRequest,
					From:    p.Self,
					To:      other,
					Content: map[string]interface{}{"value": value},
				}
				err = client.Handle(m)
				if err != nil {
					log.Warn().Err(err).Msg("handle prepare response with err")
				}
			}
		}
	}
	return nil
}

func (p *Paxos) handleVoteRequest(m Message) error {
	var err error

	compare := m.Term.Compare(p.Term)
	if compare == -1 {
		m = Message{
			Term: p.Term,
			Type: VoteResponse,
			From: m.To,
			To:   m.From,
			Ok:   false,
		}
		goto end
	}

	if compare == 1 {
		err = p.putTerm(m.Term)
		if err != nil {
			return err
		}
	}

	err = p.putVote(Vote{
		Term:  p.Term,
		Value: m.Content["value"].(Value),
	})
	if err != nil {
		return err
	}

	m = Message{
		Term: p.Term,
		Type: VoteResponse,
		From: m.To,
		To:   m.From,
		Ok:   true,
	}

end:
	client := p.Others[m.To]
	err = client.Handle(m)
	if err != nil {
		log.Warn().Err(err).Msg("handle vote request with error")
	}
	return nil
}

func (p *Paxos) handleVoteResponse(m Message) error {
	var err error

	compare := m.Term.Compare(p.Term)
	if compare == -1 {
		return nil
	}

	if compare == 1 {
		return p.putTerm(m.Term)
	}

	if compare == 0 && p.TermRound == VoteRound && m.Ok {
		p.TermAcks++
		if p.TermAcks == p.VoteNum {
			p.TermRound = NullRound
			p.TermAcks = 0

			err = p.putValue(p.Vote.Value)
			if err != nil {
				return err
			}

			for other, client := range p.Others {
				m = Message{
					Term:    p.Term,
					Type:    Learn,
					From:    p.Self,
					To:      other,
					Content: map[string]interface{}{"value": p.Value},
				}
				err = client.Handle(m)
				if err != nil {
					log.Warn().Err(err).Msg("handle vote response with error")
				}
			}
		}
	}
	return nil
}

func (p *Paxos) handleLearn(m Message) error {
	return p.putValue(m.Content["value"].(Value))
}

func (p *Paxos) Stop() {
	p.Storage.Stop()
}

package grpc

import (
	"context"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/gaoxinge/basic-paxos/pkg/core"
	"github.com/gaoxinge/basic-paxos/pkg/util/log"
)

func TermCoreToGrpc(coreTerm core.Term) *Term {
	term := Term{
		Id: coreTerm.Id,
		No: int64(coreTerm.No),
	}
	return &term
}

func TermGrpcToCore(grpcTerm *Term) core.Term {
	term := core.Term{
		Id: grpcTerm.Id,
		No: int(grpcTerm.No),
	}
	return term
}

func VoteCoreToGrpc(coreVote core.Vote) *Vote {
	vote := Vote{
		Term:  TermCoreToGrpc(coreVote.Term),
		Value: string(coreVote.Value),
	}
	return &vote
}

func VoteGrpcToCore(grpcVote *Vote) core.Vote {
	vote := core.Vote{
		Term:  TermGrpcToCore(grpcVote.Term),
		Value: core.Value(grpcVote.Value),
	}
	return vote
}

func MessageCoreToGrpc(coreMessage core.Message) *Message {
	var content map[string]*Message_Content
	if coreMessage.Content != nil {
		content = make(map[string]*Message_Content, 1)

		result, ok := coreMessage.Content["vote"]
		if ok {
			vote := result.(core.Vote)
			content["vote"] = &Message_Content{Result: &Message_Content_Vote{VoteCoreToGrpc(vote)}}
		}

		result, ok = coreMessage.Content["value"]
		if ok {
			value := result.(core.Value)
			content["value"] = &Message_Content{Result: &Message_Content_Value{string(value)}}
		}
	}

	message := Message{
		Term:    TermCoreToGrpc(coreMessage.Term),
		Type:    int64(coreMessage.Type),
		From:    coreMessage.From,
		To:      coreMessage.To,
		Ok:      coreMessage.Ok,
		Content: content,
	}
	return &message
}

func MessageGrpcToCore(grpcMessage *Message) core.Message {
	var content map[string]interface{}
	if grpcMessage.Content != nil {
		content = make(map[string]interface{}, 1)

		result, ok := grpcMessage.Content["vote"]
		if ok {
			vote := result.GetVote()
			content["vote"] = VoteGrpcToCore(vote)
		}

		result, ok = grpcMessage.Content["value"]
		if ok {
			value := result.GetValue()
			content["value"] = core.Value(value)
		}
	}

	message := core.Message{
		Term:    TermGrpcToCore(grpcMessage.Term),
		Type:    int(grpcMessage.Type),
		From:    grpcMessage.From,
		To:      grpcMessage.To,
		Ok:      grpcMessage.Ok,
		Content: content,
	}
	return message
}

type Client struct {
	Address string
}

func NewClient(address string) *Client {
	client := Client{
		Address: address,
	}
	return &client
}

func (client *Client) Handle(message core.Message) error {
	cc, err := grpc.Dial(client.Address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer cc.Close()

	c := NewServiceClient(cc)
	_, err = c.Handle(context.Background(), MessageCoreToGrpc(message))
	return err
}

type Server struct {
	Paxos     *core.Paxos
	PaxosLock sync.Mutex

	Ticker     *time.Ticker
	TickerDone chan struct{}

	GrpcListener net.Listener
	GrpcServer   *grpc.Server

	wg sync.WaitGroup
}

func NewServer(self string, others []string, path string, prepareNum int, voteNum int, first string, timeout int) (*Server, error) {
	server := &Server{}

	others0 := make(map[string]core.Client, len(others))
	for _, other := range others {
		others0[other] = NewClient(other)
	}
	first0 := core.Value(first)
	paxos, err := core.NewPaxos(self, others0, path, prepareNum, voteNum, first0)
	if err != nil {
		return nil, err
	}

	grpcListener, err := net.Listen("tcp", self)
	if err != nil {
		return nil, err
	}
	grpcServer := grpc.NewServer()
	RegisterServiceServer(grpcServer, server)

	server.Paxos = paxos
	server.Ticker = time.NewTicker(time.Duration(timeout) * time.Second)
	server.TickerDone = make(chan struct{})
	server.GrpcListener = grpcListener
	server.GrpcServer = grpcServer
	return server, nil
}

func (server *Server) Handle(ctx context.Context, message *Message) (*Empty, error) {
	coreMessage := MessageGrpcToCore(message)

	server.wg.Add(1)
	go func() {
		defer server.wg.Done()
		server.PaxosLock.Lock()
		log.Debug().Msgf("server handle %v", coreMessage)
		err := server.Paxos.Handle(coreMessage)
		if err != nil {
			log.Error().Err(err).Msg("server handle with error")
		}
		log.Debug().Msgf("server paxos with value %s", server.Paxos.Value)
		server.PaxosLock.Unlock()
	}()

	return &Empty{}, nil
}

func (server *Server) Run() {
	server.wg.Add(1)
	go func() {
		defer server.wg.Done()
		for {
			select {
			case <-server.Ticker.C:
				message := core.Message{Type: core.Timeout}

				server.PaxosLock.Lock()
				log.Debug().Msgf("server handle %v", message)
				err := server.Paxos.Handle(message)
				if err != nil {
					log.Error().Err(err).Msg("server handle with error")
				}
				log.Debug().Msgf("server paxos with value %s", server.Paxos.Value)
				server.PaxosLock.Unlock()
			case <-server.TickerDone:
				return
			}
		}
	}()

	server.wg.Add(1)
	go func() {
		defer server.wg.Done()
		err := server.GrpcServer.Serve(server.GrpcListener)
		if err != nil {
			log.Warn().Err(err).Msg("server close with error")
		}
	}()
}

func (server *Server) Stop() {
	close(server.TickerDone)
	server.GrpcServer.GracefulStop()
	server.wg.Wait()
	server.Ticker.Stop()
}

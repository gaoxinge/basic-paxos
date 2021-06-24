package core

import (
	"math/rand"
	"time"

	"github.com/gaoxinge/basic-paxos/pkg/util/log"
)

const (
	orderProbability     = 0.2
	lostProbability      = 0.2
	duplicateProbability = 0.2
	restartProbability   = 0.1
	timeoutProbability   = 0.1
)

var (
	globalRand = rand.New(rand.NewSource(time.Now().Unix()))
)

type TestClient struct {
	Messages map[string][]Message
}

func NewTestClient(messages map[string][]Message) *TestClient {
	return &TestClient{Messages: messages}
}

func (testClient *TestClient) Handle(message Message) error {
	probability := globalRand.Float64()

	if probability < orderProbability {
		if len(testClient.Messages[message.To]) > 0 {
			i := globalRand.Intn(len(testClient.Messages[message.To]))
			messages := testClient.Messages[message.To][:i]
			messages = append(messages, message)
			messages = append(messages, testClient.Messages[message.To][i:]...)
			testClient.Messages[message.To] = messages
		}
		return nil
	}

	if probability < orderProbability+lostProbability {
		return nil
	}

	if probability < orderProbability+lostProbability+duplicateProbability {
		testClient.Messages[message.To] = append(testClient.Messages[message.To], message)
		testClient.Messages[message.To] = append(testClient.Messages[message.To], message)
		return nil
	}

	testClient.Messages[message.To] = append(testClient.Messages[message.To], message)
	return nil
}

type Config struct {
	Self       string
	Others     map[string]Client
	Path       string
	PrepareNum int
	VoteNum    int
	First      Value
}

func NewConfig(self string, others map[string]Client, path string, prepareNum int, voteNum int, first Value) *Config {
	return &Config{
		Self:       self,
		Others:     others,
		Path:       path,
		PrepareNum: prepareNum,
		VoteNum:    voteNum,
		First:      first,
	}
}

type TestServer struct {
	Ids      []string
	Messages map[string][]Message
	Configs  map[string]*Config
	Paxos    map[string]*Paxos
	Value    Value
}

func NewTestServer(ids []string, paths []string, firsts []string, prepareNum int, voteNum int) (*TestServer, error) {
	messages := make(map[string][]Message, len(ids))
	for _, self := range ids {
		messages[self] = make([]Message, 0)
	}

	configs := make(map[string]*Config, len(ids))
	for i, self := range ids {
		others := make(map[string]Client)
		for _, other := range ids {
			if other != self {
				others[other] = NewTestClient(messages)
			}
		}
		configs[self] = NewConfig(self, others, paths[i], prepareNum, voteNum, Value(firsts[i]))
	}

	paxos := make(map[string]*Paxos, len(ids))
	for _, self := range ids {
		config := configs[self]
		p, err := NewPaxos(config.Self, config.Others, config.Path, config.PrepareNum, config.VoteNum, config.First)
		if err != nil {
			return nil, err
		}
		paxos[self] = p
	}

	return &TestServer{
		Ids:      ids,
		Messages: messages,
		Configs:  configs,
		Paxos:    paxos,
	}, nil
}

func (testServer *TestServer) Handle() error {
	probability := globalRand.Float64()

	if probability < restartProbability {
		i := globalRand.Intn(len(testServer.Ids))
		id := testServer.Ids[i]

		log.Info().Msgf("[HANDLE] paxos %s handle restart", id)

		config := testServer.Configs[id]
		p := testServer.Paxos[id]
		p.Stop()
		p, err := NewPaxos(config.Self, config.Others, config.Path, config.PrepareNum, config.VoteNum, config.First)
		if err != nil {
			return err
		}
		testServer.Paxos[id] = p
		return nil
	}

	if probability < restartProbability+timeoutProbability {
		i := globalRand.Intn(len(testServer.Ids))
		id := testServer.Ids[i]

		log.Info().Msgf("[HANDLE] paxos %s handle timeout", id)

		p := testServer.Paxos[id]
		return p.Handle(Message{Type: Timeout})
	}

	ids := make([]string, 0, len(testServer.Ids))
	for _, id := range testServer.Ids {
		if len(testServer.Messages[id]) > 0 {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return nil
	}
	i := globalRand.Intn(len(ids))
	id := ids[i]

	p := testServer.Paxos[id]
	message := testServer.Messages[id][0]
	testServer.Messages[id] = testServer.Messages[id][1:]

	log.Info().Msgf("[HANDLE] paxos %s handle message %v", id, message)
	return p.Handle(message)
}

func (testServer *TestServer) Check() bool {
	for _, id := range testServer.Ids {
		value := testServer.Paxos[id].Value
		if value.IsNull() {
			continue
		}

		log.Info().Msgf("[CHECK] paxos %s with value %v", id, value)
		if testServer.Value.IsNull() {
			testServer.Value = value
			continue
		} else {
			if testServer.Value != value {
				return false
			}
		}
	}
	return true
}

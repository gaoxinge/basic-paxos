package http

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/gaoxinge/basic-paxos/pkg/core"
	"github.com/gaoxinge/basic-paxos/pkg/util/log"
)

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
	bs, err := jsoniter.Marshal(message)
	if err != nil {
		return err
	}

	response, err := http.Post(fmt.Sprintf("http://%s/handle", client.Address), "application/json", bytes.NewReader(bs))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.Errorf("client get http status code %d", response.StatusCode)
	}
	return nil
}

type Server struct {
	Paxos     *core.Paxos
	PaxosLock sync.Mutex

	Ticker     *time.Ticker
	TickerDone chan struct{}

	HttpServer *http.Server

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

	gin.SetMode(gin.DebugMode)
	router := gin.New()
	router.POST("/handle", server.Handle)
	httpServer := http.Server{
		Addr:    self,
		Handler: router,
	}

	server.Paxos = paxos
	server.Ticker = time.NewTicker(time.Duration(timeout) * time.Second)
	server.TickerDone = make(chan struct{})
	server.HttpServer = &httpServer
	return server, nil
}

func (server *Server) Handle(ctx *gin.Context) {
	var message core.Message
	err := jsoniter.NewDecoder(ctx.Request.Body).Decode(&message)
	if err != nil {
		log.Error().Err(err).Msg("server handle with error")
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
	content := message.Content
	if content != nil {
		result, ok := content["vote"]
		if ok {
			vote := result.(map[string]interface{})
			term := vote["term"].(map[string]interface{})
			value := vote["value"].(string)
			term0 := core.Term{Id: term["id"].(string), No: int(term["no"].(float64))}
			value0 := core.Value(value)
			vote0 := core.Vote{Term: term0, Value: value0}
			message.Content = map[string]interface{}{"vote": vote0}
		}

		result, ok = content["value"]
		if ok {
			value := result.(string)
			value0 := core.Value(value)
			message.Content = map[string]interface{}{"value": value0}
		}
	}
	ctx.JSON(http.StatusOK, nil)

	server.wg.Add(1)
	go func() {
		defer server.wg.Done()
		server.PaxosLock.Lock()
		log.Debug().Msgf("server handle %v", message)
		err = server.Paxos.Handle(message)
		if err != nil {
			log.Error().Err(err).Msg("server handle with error")
		}
		log.Debug().Msgf("server paxos with value %s", server.Paxos.Value)
		server.PaxosLock.Unlock()
	}()
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
		err := server.HttpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Warn().Err(err).Msg("server close with error")
		}
	}()
}

func (server *Server) Stop() {
	close(server.TickerDone)
	err := server.HttpServer.Shutdown(context.TODO())
	if err != nil && err != http.ErrServerClosed {
		log.Warn().Err(err).Msg("server stop with error")
	}

	server.wg.Wait()
	server.Ticker.Stop()
}

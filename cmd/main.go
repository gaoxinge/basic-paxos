package main

import (
	"flag"
	"strings"

	"github.com/gaoxinge/basic-paxos/pkg/core"
	"github.com/gaoxinge/basic-paxos/pkg/grpc"
	"github.com/gaoxinge/basic-paxos/pkg/http"
	"github.com/gaoxinge/basic-paxos/pkg/util/log"
)

func main() {
	kindPtr := flag.String("kind", "", "kind")
	selfPtr := flag.String("self", "", "address of self")
	othersPtr := flag.String("others", "", "address of other")
	pathPtr := flag.String("path", "", "path of storage")
	prepareNumPtr := flag.Int("prepareNum", 0, "prepare of num")
	voteNumPtr := flag.Int("voteNum", 0, "vote of num")
	firstPtr := flag.String("first", "", "value")
	timeoutPtr := flag.Int("timeout", 0, "timeout")
	flag.Parse()

	kind := *kindPtr
	self := *selfPtr
	others := strings.Split(*othersPtr, ",")
	path := *pathPtr
	prepareNum := *prepareNumPtr
	voteNum := *voteNumPtr
	first := *firstPtr
	timeout := *timeoutPtr
	log.Debug().Msgf("kind: %s, self: %s, others: %v, path: %s, prepareNum: %d, voteNum: %d, first: %s, timeout: %d",
		kind, self, others, path, prepareNum, voteNum, first, timeout)

	var (
		server core.Server
		err    error
	)
	switch kind {
	case "http":
		server, err = http.NewServer(self, others, path, prepareNum, voteNum, first, timeout)
		if err != nil {
			log.Error().Err(err).Msg("new http server with error")
			return
		}
	case "grpc":
		server, err = grpc.NewServer(self, others, path, prepareNum, voteNum, first, timeout)
		if err != nil {
			log.Error().Err(err).Msg("new http server with error")
			return
		}
	default:
		log.Error().Msgf("unknown kind %s", kind)
		return
	}

	server.Run()
	select {}
}

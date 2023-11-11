package node

import (
	"context"
	"flag"
	"fmt"
	"game-lb/config"
	"game-lb/request"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/share"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"
)

var gameNodeServer []*string

var nodeXClient client.XClient

func Init() {
	initGameHost()
	initXClient()
}

func initGameHost() {
	for i := 0; i < config.Conf.GamePodNumb; i++ {
		host := fmt.Sprintf("%s-%d.%s.%s.svc.cluster.local", config.Conf.GamePodName, i, config.Conf.GamePodService, config.Conf.GamePodNamespace)
		if config.RunMode == "debug" {
			host = "127.0.0.1"
		}
		clientAddress := flag.String(fmt.Sprintf("%s-%d", config.Conf.GamePodName, i), fmt.Sprintf("tcp@%s:8888", host), fmt.Sprintf("%s-%d", config.Conf.GamePodName, i))
		gameNodeServer = append(gameNodeServer, clientAddress)
	}
}

func initXClient() {
	flag.Parse()
	var clients []*client.KVPair
	for i := range gameNodeServer {
		clients = append(clients, &client.KVPair{Key: *gameNodeServer[i]})
	}
	d, _ := client.NewMultipleServersDiscovery(clients)
	nodeXClient = client.NewXClient("Arith", client.Failfast, client.SelectByUser, d, client.Option{
		Retries:             3,
		TimeToDisallow:      time.Minute,
		RPCPath:             share.DefaultRPCPath,
		ConnectTimeout:      time.Second,
		IdleTimeout:         5 * time.Second,
		BackupLatency:       10 * time.Millisecond,
		GenBreaker:          nil,
		SerializeType:       protocol.MsgPack,
		CompressType:        protocol.None,
		Heartbeat:           true,
		HeartbeatInterval:   time.Second,
		MaxWaitForHeartbeat: time.Second,
		TCPKeepAlivePeriod:  time.Minute,
		BidirectionalBlock:  false,
	})
	nodeXClient.SetSelector(&gameIdSelector{})
	cleanupHook()
}

type gameIdSelector struct {
	servers []string
}

func (s *gameIdSelector) Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string {
	var ss = s.servers
	if len(ss) == 0 {
		return ""
	}
	var gameId uint
	switch serviceMethod {
	case "InitInfo":
		gameId = args.(request.InitRequest).GameId
	case "Info":
		gameId = args.(request.InfoRequest).GameId
	case "Enter":
		gameId = args.(request.EnterRequest).GameId
	case "Move":
		gameId = args.(request.MoveRequest).GameId
	case "Pass":
		gameId = args.(request.PassRequest).GameId
	case "Resign":
		gameId = args.(request.ResignRequest).GameId
	case "AreaScore":
		gameId = args.(request.InfoRequest).GameId
	case "ApplyAreaScore":
		gameId = args.(request.InfoRequest).GameId
	case "RejectAreaScore":
		gameId = args.(request.InfoRequest).GameId
	case "AgreeAreaScore":
		gameId = args.(request.InfoRequest).GameId
	case "SGFInfo":
		gameId = args.(request.InfoRequest).GameId
	case "CallEnd":
		gameId = args.(request.CallEndRequest).GameId
	case "ApplySummation":
		gameId = args.(request.InfoRequest).GameId
	case "RejectSummation":
		gameId = args.(request.InfoRequest).GameId
	case "AgreeSummation":
		gameId = args.(request.InfoRequest).GameId
	case "CanPlay":
		gameId = args.(request.InfoRequest).GameId
	case "InnerAreaScore":
		gameId = args.(request.OwnershipRequest).GameId
	case "Ownership":
		gameId = args.(request.OwnershipRequest).GameId
	case "ForceReloadSGF":
		gameId = args.(request.ForceReloadSgfRequest).GameId
	default:
		gameId = 0
	}
	gameIdStr := fmt.Sprintf("%d", gameId)
	var target string
	hash := 0
	for i := 0; i < len(gameIdStr); i++ {
		hash += int(gameIdStr[i])
	}
	index := hash % len(ss)
	if index >= 0 && index < len(ss) {
		target = ss[index]
	} else {
		target = ss[0]
	}
	return target
}

func (s *gameIdSelector) UpdateServer(servers map[string]string) {
	var ss = make([]string, 0, len(servers))
	for k := range servers {
		ss = append(ss, k)
	}

	sort.Slice(ss, func(i, j int) bool {
		return strings.Compare(ss[i], ss[j]) <= 0
	})
	s.servers = ss
}

func cleanupHook() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		_ = nodeXClient.Close()
		os.Exit(0)
	}()
}

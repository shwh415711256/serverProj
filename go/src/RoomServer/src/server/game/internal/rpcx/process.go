package rpcx

import (
	"golang.org/x/net/context"
	"Common/model"
	"server/game/internal/room"
	"github.com/cihub/seelog"
	"github.com/smallnest/rpcx/serverplugin"
	"time"
	"github.com/smallnest/rpcx/server"
	"server/conf"
	"github.com/rcrowley/go-metrics"
	"log"
)

type RpcProcess struct{}

var (
	MRcpProcess *RpcProcess
)

// 注册到服务发现
func addRegistryPlugin(server *server.Server) {
	etcd := &serverplugin.EtcdRegisterPlugin{
		ServiceAddress:"tcp@127.0.0.1" + conf.Server.RpcxPort,
		EtcdServers:[]string{conf.Server.EtcdAddr},
		BasePath:"/zq/rpcx",
		Metrics:metrics.NewRegistry(),
		UpdateInterval:time.Minute,
	}
	err := etcd.Start()
	if err != nil {
		log.Fatal(err)
	}
	server.Plugins.Add(etcd)
}

func Init() {
	// 先用rpcx点对点吧，定死一个roomserver,之后再去改
	s := server.NewServer()
	addRegistryPlugin(s)
	s.RegisterName(conf.Server.RoomServerPath, new(RpcProcess), "")
	go s.Serve("tcp", "127.0.0.1"+conf.Server.RpcxPort)
}

func (p *RpcProcess) RpcCreateRoom (ctx context.Context, req *model.CreateRoomRequest, reply *model.CreateRoomResponse) error {
	err := room.CreateOneRoom(req, reply)
	if err != nil {
		seelog.Errorf("create one room error[%v], err")
		return err
	}
	return nil
}
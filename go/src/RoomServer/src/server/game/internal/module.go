package internal

import (
	"github.com/name5566/leaf/module"
	"server/base"
	"server/game/internal/room"
	"server/game/internal/rpcclient"
	"server/game/internal/rpcx"
	"server/game/internal/timetick"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
	rpcx.Init()
	room.Init()
	rpcclient.Init()
	timetick.Init()
}

func (m *Module) OnDestroy() {
	rpcclient.Close()
}

package timetick

import (
	"github.com/name5566/leaf/timer"
	"server/game/internal/room"
)

var (
	OneSecTimeTick *timer.Dispatcher
)

func Init() {
	OneSecTimeTick = timer.NewDispatcher(5)
	// 每秒
	cronExpr, err := timer.NewCronExpr("* * * * * *")
	if err != nil {
		return
	}
	OneSecTimeTick.CronFunc(cronExpr, func() {
		room.UpdateOneSec()
	})
	go func() {
		for {
			(<-OneSecTimeTick.ChanTimer).Cb()
		}
	}()
}

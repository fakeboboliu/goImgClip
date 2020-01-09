package main

import (
	"fmt"
	"github.com/popu125/goImgClip/log"
	"github.com/popu125/goImgClip/manager"
	"github.com/popu125/goImgClip/notify"
	T "github.com/popu125/goImgClip/targets"
	"github.com/popu125/goImgClip/tray"
	"os"

	"github.com/getlantern/systray"
	"github.com/rs/zerolog"
)

func main() {
	// Set a func to recover and show error message
	defer func() {
		if e := recover(); e != nil {
			notify.Error(fmt.Sprint(e))
			os.Exit(1)
		}
	}()

	t := tray.NewTray()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	conf := LoadConfig("config.yml")
	log.InitLogger(conf.LogFile)
	l := log.GetLogger("main")
	check := log.GenCheck(l)

	targets := make(map[string]T.Target, len(conf.Targets))
	l.Debug().Int("amount", len(conf.Targets)).Msg("Loading targets")

	for _, tc := range conf.Targets {
		tmpTarget := T.NewTarget(tc.Target)
		check(tmpTarget.Configure(tc.TargetConfig))
		targets[tc.Name] = tmpTarget

		l.Debug().Str("name", tc.Name).Msg("Target loaded")
	}

	m := manager.NewManager(targets, t)
	m.ListenOnKey(conf.HotKey)
	go systray.Run(m.OnTrayReady, m.OnTrayExit)

	m.WaitExit()
}

package manager

import (
	"bytes"
	"fmt"
	"github.com/popu125/goImgClip/clipboard"
	"github.com/popu125/goImgClip/hotkey"
	"github.com/popu125/goImgClip/log"
	"github.com/popu125/goImgClip/notify"
	T "github.com/popu125/goImgClip/targets"
	"github.com/popu125/goImgClip/tray"
	"github.com/rs/zerolog"
	"image/jpeg"
)

type manager struct {
	targets map[string]T.Target
	names   []string
	curT    T.Target

	hkStr string
	hkLis hotkey.Listener
	t     *tray.Tray
	l     zerolog.Logger

	trayReady chan bool
	trayExit  chan bool
}

func NewManager(targets map[string]T.Target, tray *tray.Tray) *manager {
	names := make([]string, len(targets))
	i := 0
	for name := range targets {
		names[i] = name
		i++
	}

	return &manager{targets: targets,
		names:     names,
		hkLis:     hotkey.KeyListener{}.New(),
		t:         tray,
		curT:      nil,
		l:         log.GetLogger("mgr"),
		trayReady: make(chan bool),
		trayExit:  make(chan bool),
	}
}

func (m *manager) ListenOnTray() {
	for name, target := range m.targets {
		if m.curT == nil {
			m.curT = target
		}
		m.t.SetItem(name, fmt.Sprint("Target: ", name))
		m.t.AddListener(name, m.genTrayListener(name))
	}
}

func (m *manager) ListenOnKey(keys string) {
	m.hkStr = keys
	err := m.hkLis.Listen("main", keys, m.ProcessUpload)
	if err != nil {
		m.l.Warn().Str("err", err.Error()).Msg("ListenOnKey failed")
		notify.Error("ListenOnKey failed")
		return
	}
}

func (m *manager) genTrayListener(name string) func(t *tray.Tray) {
	return func(t *tray.Tray) {
		// uncheck all and check it
		for _, tmp := range m.names {
			if item, ok := t.GetItem(tmp); ok {
				item.Uncheck()
			}
		}
		if item, ok := t.GetItem(name); ok {
			item.Check()
		}

		// and let hotkey listen it
		m.setTarget(name)
	}
}

func (m *manager) setTarget(name string) {
	m.curT = m.targets[name]
}

func (m *manager) ProcessUpload() {
	notify.Action("Working")

	rgba, err := clipboard.GetImageFromClipBoard()
	if err != nil {
		m.l.Warn().Str("err", err.Error()).Msg("GetImageFromClipBoard failed")
		notify.Error("GetImageFromClipBoard failed")
		return
	}
	if rgba == nil {
		notify.Error("There're no image in clipboard")
		return
	}

	img := rgba.SubImage(rgba.Rect)
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 100})
	if err != nil {
		m.l.Warn().Str("err", err.Error()).Msg("Encode image from clipboard failed")
		notify.Error("Encode image from clipboard failed")
		return
	}

	url, err := m.curT.Upload(buf.Bytes())
	if err != nil {
		m.l.Warn().Str("err", err.Error()).Msg("Upload failed")
		notify.Error("Upload failed")
		return
	}

	err = clipboard.SetTextToClipboard(url)
	if err != nil {
		m.l.Warn().Str("err", err.Error()).Msg("Set clipboard failed")
		notify.Error("Set clipboard failed")
		return
	}

	notify.Success(m.curT.Name())
}

func (m *manager) WaitExit() {
	<-m.trayExit
}

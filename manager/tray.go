package manager

import (
	"fmt"
	"github.com/popu125/goImgClip/tray"

	"github.com/getlantern/systray"
)

func (m *manager) OnTrayReady() {
	systray.SetIcon(icon)
	systray.SetTitle("goImgClip")
	systray.SetTooltip("goImgClip - A tool to upload image from clipboard")
	m.t.SetItem("hotkey", fmt.Sprint("Hotkey: ", m.hkStr)).Disable()
	m.t.AddSeparator()
	m.ListenOnTray()
	m.t.AddSeparator()
	m.t.SetItem("quit", "Quit")
	m.t.AddListener("quit", func(t *tray.Tray) {
		systray.Quit()
	})
	m.trayReady <- true
}

func (m *manager) OnTrayExit() {
	m.trayExit <- true
}

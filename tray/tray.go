package tray

import (
	"sync"
	"time"

	"github.com/getlantern/systray"
)

type menuItem struct {
	*systray.MenuItem
	listened bool
}

type itemMap map[string]menuItem
type Tray struct {
	items itemMap
	wg    sync.WaitGroup // WaitGroup for listeners
	icon  []byte
}

func NewTray() *Tray {
	return &Tray{items: make(itemMap)}
}

func (t *Tray) SetItem(name string, title string) menuItem {
	if v, ok := t.GetItem(name); ok {
		return v
	}
	item := systray.AddMenuItem(title, "")
	t.items[name] = menuItem{item, false}
	return t.items[name]
}

func (t *Tray) GetItem(name string) (menuItem, bool) {
	v, ok := t.items[name]
	return v, ok
}

func (t *Tray) AddListener(name string, listener func(*Tray)) bool {
	if v, ok := t.GetItem(name); ok && !v.listened {
		go func() {
			for range v.ClickedCh {
				listener(t)
			}
		}()
		return true
	}
	return false
}

func (t *Tray) AddSeparator() {
	systray.AddSeparator()
}

func (t *Tray) ShowIcon(data []byte, duration time.Duration) {
	// todo
}

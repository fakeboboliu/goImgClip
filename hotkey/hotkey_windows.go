package hotkey

import (
	"errors"
	"github.com/MakeNowJust/hotkey"
	"strconv"
	"strings"
)

var strToModsMap = map[string]hotkey.Modifier{
	"ctrl":  hotkey.Ctrl,
	"alt":   hotkey.Alt,
	"shift": hotkey.Shift,
	"win":   hotkey.Win,
}

type keyGroup struct {
	mod       hotkey.Modifier
	vk        uint32
	listening bool
	id        hotkey.Id
}

func (g *keyGroup) setKey(str string) error {
	keys := strings.Split(str, "+")
	g.mod, g.vk = hotkey.Modifier(0), uint32(0)

	for i, key := range keys {
		// parse modifiers
		k := strings.ToLower(strings.TrimSpace(key))
		switch k {
		case "ctrl", "alt", "shift", "win":
			g.mod += strToModsMap[k]
			continue
		}

		if i+1 == len(keys) { // parse vk here
			if strings.HasPrefix(k, "f") && len(k) > 1 {
				num, err := strconv.Atoi(k[1:])
				if err == nil && num > 0 {
					g.vk = hotkey.F1 + uint32(num) - 1
				}
			} else {
				g.vk = uint32(strings.ToUpper(k)[0])
			}
		}
	}

	if g.mod == 0 || g.vk == 0 {
		return errors.New("there should be both type of key in hotkey")
	}

	return nil
}

type KeyListener struct {
	names map[string]*keyGroup
	hkey  *hotkey.Manager
}

func (l *KeyListener) getKey(name string) *keyGroup {
	if v, ok := l.names[name]; ok {
		return v
	}

	kg := &keyGroup{listening: false}
	l.names[name] = kg
	return kg
}

func (l *KeyListener) Listen(name string, keys string, cb func()) error {
	if l.hkey == nil {
		return errors.New("instance of Listener should be init by New()")
	}

	var err error
	kg := l.getKey(name)

	if kg.listening {
		err = l.UnListen(name)
		if err != nil {
			return err
		}
	}

	err = kg.setKey(keys)
	if err != nil {
		return err
	}

	kg.id, err = l.hkey.Register(kg.mod, kg.vk, cb)
	return err
}

func (l *KeyListener) UnListen(name string) error {
	kg := l.getKey(name)
	if kg.listening {
		l.hkey.Unregister(kg.id)
	}
	return nil
}

func (l KeyListener) New() Listener {
	return &KeyListener{hkey: hotkey.New(), names: map[string]*keyGroup{}}
}

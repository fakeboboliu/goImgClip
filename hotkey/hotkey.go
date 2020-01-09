package hotkey

type Listener interface {
	New() Listener
	Listen(name string, keys string, cb func()) error
	UnListen(name string) error
}

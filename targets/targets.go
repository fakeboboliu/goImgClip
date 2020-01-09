package targets

import "gopkg.in/yaml.v3"

type Target interface {
	New(name string) Target
	Name() string
	Upload(img []byte) (string, error)
	Configure(raw yaml.Node) error
}

var nameToTarget = map[string]Target{}

func registerTarget(name string, target Target) {
	nameToTarget[name] = target
}

func NewTarget(targetName string) Target {
	if v, ok := nameToTarget[targetName]; ok {
		return v.New(targetName)
	}
	return nil
}

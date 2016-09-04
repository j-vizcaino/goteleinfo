package exporters

import (
	"goteleinfo/teleinfo"
)

type Factory func() (teleinfo.Exporter, error)

var registry map[string]Factory

func init() {
	registry = map[string]Factory{}
}

func Register(name string, f Factory) {
	registry[name] = f
}

func List() []string {
	res := make([]string, len(registry))
	idx := 0
	for name := range registry {
		res[idx] = name
		idx++
	}
	return res
}

func FindByName(name string) Factory {
	return registry[name]
}

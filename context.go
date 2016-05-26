package summer

import "reflect"

type Stone interface{}

type Plugin interface {
	Look(path string) reflect.Value
	Prefix() string
	ZIndex() int
}
type PluginWorkTime int

const (
	BeforeInit PluginWorkTime = iota
	AfterInit
)

type Basket interface {
	Add(name string, stone Stone)
	Put(stone Stone)
	Stone(name string, t reflect.Type) (stone Stone)
	NameStone(name string) (stone Stone)
	PluginRegister(Plugin,PluginWorkTime)
	Start()
	ShutDown()
}
type Init interface {
	Init()
}
type Ready interface {
	Ready()
}
type Destroy interface {
	Destroy()
}
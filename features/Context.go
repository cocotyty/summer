package features

import "reflect"

type Stone interface{}

type Plugin interface {
	Look(path string) reflect.Value
	Prefix() string
	ZIndex() int
}


type Basket interface {
	Add(name string, stone Stone)
	Put(stone Stone)
	Stone(name string, t reflect.Type) (stone Stone)
	NameStone(name string) (stone Stone)
	PluginRegister(Plugin)
	Start()
	ShutDown()
}
type Init interface {
	Init()
}
type Destroy interface {
	Destroy()
}
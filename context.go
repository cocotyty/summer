package summer

import "reflect"

type Stone interface{}


type Basket interface {
	Add(name string, stone Stone)
	Put(stone Stone)
	Stone(name string, t reflect.Type) (stone Stone)
	NameStone(name string) (stone Stone)
	NameHolder(name string) (holder *Holder)
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

package features

import "reflect"

type Stone interface{}

type Plugin interface {
	Look(path string) interface{}
	Prefix() string
	ZIndex() int
}

type PathFactory interface {
	Name()
	Stone(path string, value reflect.Value) Stone
}

type Basket interface {
	Add(name string, stone Stone)
	Put(stone Stone)
	Stone(name string, t reflect.Type) (stone Stone)
	PathFactoryRegister(PathFactory)
	Start()
	ShutDown()
}
type Init interface {
	Init()
}
type Destroy interface {
	Destroy()
}
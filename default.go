package summer

import "reflect"

var defaultBasket = NewBasket()

// add a stone to the default basket with a name
func Add(name string, stone Stone, value ...interface{}) {
	if len(value) == 0 {
		defaultBasket.Add(name, stone)
	} else {
		defaultBasket.AddWithValue(name, stone, value[0])
	}
}

// put a stone into the default basket
func Put(stone Stone, value ...interface{}) {
	if len(value) == 0 {
		defaultBasket.Put(stone)
	} else {
		defaultBasket.PutWithValue(stone, value[0])
	}
}

// get a stone with the name and the type
func GetStone(name string, t reflect.Type) (stone Stone) {
	return defaultBasket.GetStone(name, t)
}

// get a tone with the name
func GetStoneWithName(name string) (stone Stone) {
	return defaultBasket.GetStoneWithName(name)
}

// register a plugin to basket
func PluginRegister(p Plugin, pt PluginWorkTime) {
	defaultBasket.PluginRegister(p, pt)
}
func Start() {
	defaultBasket.Start()
}
func Strict() {
	defaultBasket.Strict()
}
func ShutDown() {
	defaultBasket.ShutDown()
}

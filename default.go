package summer

import "reflect"

var defaultBasket = NewBasket()
// add a stone to the default basket with a name
func Add(name string, stone Stone) {
	defaultBasket.Add(name, stone)
}
// put a stone into the default basket
func Put(stone Stone) {
	defaultBasket.Put(stone)
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
func ShutDown() {
	defaultBasket.ShutDown()
}

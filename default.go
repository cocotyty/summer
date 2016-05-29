package summer

import "reflect"

var defaultBasket = NewBasket()

func Add(name string, stone Stone) {
	defaultBasket.Add(name, stone)
}
func Put(stone Stone) {
	defaultBasket.Put(stone)
}
func GetStone(name string, t reflect.Type) (stone Stone) {
	return defaultBasket.Stone(name, t)

}
func NameStone(name string) (stone Stone) {
	return defaultBasket.NameStone(name)
}
func PluginRegister(p Plugin, pt PluginWorkTime) {
	defaultBasket.PluginRegister(p, pt)
}
func Start() {
	defaultBasket.Start()
}
func ShutDown() {
	defaultBasket.ShutDown()
}

package summer

import "reflect"

type Plugin interface {
	// look up the value which field wanted
	Look(Holder *Holder, path string, sf *reflect.StructField) reflect.Value
	// tell  summer the plugin prefix
	Prefix() string
	// zIndex represent the sequence of plugins
	ZIndex() int
}
type PluginWorkTime int

const (
	BeforeInit PluginWorkTime = iota
	AfterInit
)

type pluginList []Plugin

func (list pluginList) Len() int {
	return len(list)
}
func (list pluginList) Less(i, j int) bool {
	return list[i].ZIndex() < list[j].ZIndex()
}
func (list pluginList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

package summer

import "reflect"
import (
	"github.com/pelletier/go-toml"
)


type TomlPlugin struct {
	tree *toml.TomlTree
}

func (this *TomlPlugin) Look(h *Holder,path string) reflect.Value {
	return reflect.ValueOf(this.tree.Get(path))
}
func (this *TomlPlugin) Prefix() string {
	return "#"
}
func (this *TomlPlugin) ZIndex() int {
	return 0
}

package summer

import "reflect"
import (
	"github.com/pelletier/go-toml"
)

type TomlPlugin struct {
	tree *toml.TomlTree
}

func (this *TomlPlugin) Look(path string) reflect.Value {
	logger.Println(path)
	return reflect.ValueOf(this.tree.Get(path))
}
func (this *TomlPlugin) Prefix() string {
	return "#"
}
func (this *TomlPlugin) ZIndex() int {
	return 0
}

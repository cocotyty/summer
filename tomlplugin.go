package summer

import "reflect"
import (
	"github.com/pelletier/go-toml"
)


func TomlFile(path string) error {
	tree, err := toml.LoadFile(path)
	if err != nil {
		return err
	}
	defaultBasket.PluginRegister(&TomlPlugin{tree}, BeforeInit)
	return nil
}

func Toml(src string) error {
	tree, err := toml.Load(src)
	if err != nil {
		return err
	}
	defaultBasket.PluginRegister(&TomlPlugin{tree}, BeforeInit)
	return nil
}

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

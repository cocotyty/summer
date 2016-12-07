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

func (this *TomlPlugin) Look(h *Holder, path string, sf *reflect.StructField) reflect.Value {
	if sf.Type.Kind() == reflect.Slice {
		v := this.tree.Get(path)
		if sf.Type.Elem().Kind() == reflect.String {
			strs := []string{}
			if list, ok := v.([]interface{}); ok {
				for _, elm := range list {
					str, ok := elm.(string)
					if !ok {
						panic("Toml is Wrong! @" + path)
					}
					strs = append(strs, str)
				}
			}
			return reflect.ValueOf(strs)
		}
		if sf.Type.Elem().Kind() == reflect.Int {
			ints := []int{}
			if list, ok := v.([]interface{}); ok {
				for _, elm := range list {
					i, ok := elm.(int64)
					if !ok {
						panic("Toml is Wrong! @" + path)
					}
					ints = append(ints, int(i))
				}
			}
			return reflect.ValueOf(ints)
		}

		if sf.Type.Elem().Kind() == reflect.Int64 {
			ints := []int64{}
			if list, ok := v.([]interface{}); ok {
				for _, elm := range list {
					i, ok := elm.(int64)
					if !ok {
						panic("Toml is Wrong! @" + path)
					}
					ints = append(ints, i)
				}
			}
			return reflect.ValueOf(ints)
		}
	}
	return reflect.ValueOf(this.tree.Get(path))
}
func (this *TomlPlugin) Prefix() string {
	return "#"
}
func (this *TomlPlugin) ZIndex() int {
	return 0
}

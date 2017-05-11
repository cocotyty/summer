package summer

import "reflect"
import (
	"github.com/pelletier/go-toml"
)

func NewTomlPluginByFilePath(path string) (*TomlPlugin, error) {
	tree, err := toml.LoadFile(path)
	if err != nil {
		return nil, err
	}
	return &TomlPlugin{tree: tree}, nil
}
func NewTomlPluginBySource(src string) (*TomlPlugin, error) {
	tree, err := toml.Load(src)
	if err != nil {
		return nil, err
	}
	return &TomlPlugin{tree: tree}, nil
}

type TomlPlugin struct {
	tree *toml.Tree
}

func (plugin *TomlPlugin) Look(h *Holder, path string, sf *reflect.StructField) reflect.Value {
	if sf.Type.Kind() == reflect.Slice {
		value := plugin.tree.Get(path)
		if sf.Type.Elem().Kind() == reflect.String {
			stringSliceValue := []string{}
			if list, ok := value.([]interface{}); ok {
				for _, elm := range list {
					str, ok := elm.(string)
					if !ok {
						panic("Toml is Wrong! @" + path)
					}
					stringSliceValue = append(stringSliceValue, str)
				}
			}
			return reflect.ValueOf(stringSliceValue)
		}
		if sf.Type.Elem().Kind() == reflect.Int {
			intSliceValue := []int{}
			if list, ok := value.([]interface{}); ok {
				for _, elm := range list {
					i, ok := elm.(int64)
					if !ok {
						panic("Toml is Wrong! @" + path)
					}
					intSliceValue = append(intSliceValue, int(i))
				}
			}
			return reflect.ValueOf(intSliceValue)
		}
		if sf.Type.Elem().Kind() == reflect.Int64 {
			int64SliceValue := []int64{}
			if list, ok := value.([]interface{}); ok {
				for _, elm := range list {
					i, ok := elm.(int64)
					if !ok {
						panic("Toml is Wrong! @" + path)
					}
					int64SliceValue = append(int64SliceValue, i)
				}
			}
			return reflect.ValueOf(int64SliceValue)
		}
	}
	return reflect.ValueOf(plugin.tree.Get(path))
}
func (plugin *TomlPlugin) Prefix() string {
	return "#"
}
func (plugin *TomlPlugin) ZIndex() int {
	return 0
}

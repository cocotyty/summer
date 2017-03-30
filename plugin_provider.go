package summer

import (
	"reflect"
)

type ProviderPlugin struct {
}

func (plugin *ProviderPlugin) Look(holder *Holder, path string, sf *reflect.StructField) (need reflect.Value) {
	if path == "*" {
		path = sf.Name
	}
	holder.Basket.EachHolder(func(name string, holder *Holder) bool {
		if provider, ok := holder.Stone.(Provider); ok {
			stone := provider.Provide()
			if name == path {
				need = reflect.ValueOf(stone)
				return true
			}

		}
		return false
	})
	empty := reflect.Value{}
	if need == empty {
		holder.Basket.EachHolder(func(name string, holder *Holder) bool {
			if provider, ok := holder.Stone.(Provider); ok {
				stone := provider.Provide()
				if reflect.TypeOf(stone) == sf.Type || reflect.TypeOf(stone).AssignableTo(sf.Type) || reflect.TypeOf(stone).ConvertibleTo(sf.Type) {
					need = reflect.ValueOf(stone)
					return true
				}
			}
			return false
		})
	}
	if need == empty {
		panic("provider not found:" + holder.Type.PkgPath() + "/" + holder.Type.Name() + "." + sf.Name + " @." + path)
	}
	return need
}
func (plugin *ProviderPlugin) Prefix() string {
	return "@"
}

// zIndex represent the sequence of plugins
func (plugin *ProviderPlugin) ZIndex() int {
	return 3
}

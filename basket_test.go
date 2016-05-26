package summer

import (
	"github.com/cocotyty/summer/simples"
	"github.com/pelletier/go-toml"
	"reflect"
	"testing"
)

func TestBasket(t *testing.T) {
	b := NewBasket()
	tom, _ := toml.Load(`[postgres]
user = "pelletier"
password = "mypassword"`)
	b.PluginRegister(&TomlPlugin{tom}, BeforeInit)
	b.PluginRegister(&RefPlugin{b}, AfterInit)
	b.Put(&simples.A{})
	b.Put(&simples.C{})
	b.Put(&simples.D{})
	b.Start()
	b.Stone("A", reflect.TypeOf(&simples.A{})).(*simples.A).Print()
	b.ShutDown()
}

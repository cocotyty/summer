package features

import (
	"testing"
	"github.com/cocotyty/summer/simples"
	"reflect"
	"github.com/pelletier/go-toml"
)

func TestBasket(t *testing.T) {
	b := NewBasket()
	tom,_:=toml.Load(`[postgres]
user = "pelletier"
password = "mypassword"`)
	b.PluginRegister(&TomlPlugin{tom})
	b.PluginRegister(&RefPlugin{b})
	b.Put(&simples.A{})
	b.Put(&simples.C{})
	b.Put(&simples.D{})
	b.Start()
	b.Stone("A", reflect.TypeOf(&simples.A{})).(*simples.A).Print()
	b.ShutDown()
}

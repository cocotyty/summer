package summer

import (
	"log"
	"os"
	"os/signal"
	"reflect"
)

var DefaultBasket = NewBasket()

func init() {
	PluginRegister(&ProviderPlugin{}, AfterInit)
	PluginRegister(&FieldReferencePlugin{}, AfterInit)
}

func TomlFile(path string) error {
	plugin, err := NewTomlPluginByFilePath(path)
	if err != nil {
		return err
	}
	PluginRegister(plugin, BeforeInit)
	return nil
}

func Toml(src string) error {
	plugin, err := NewTomlPluginBySource(src)
	if err != nil {
		return err
	}
	PluginRegister(plugin, BeforeInit)
	return nil
}

func AddNotStrict(name string, stone Stone, value ...interface{}) {
	if len(value) == 0 {
		DefaultBasket.AddNotStrict(name, stone, nil)
	} else {
		DefaultBasket.AddNotStrict(name, stone, value[0])
	}
}
func PutNotStrict(stone Stone, value ...interface{}) {
	if len(value) == 0 {
		DefaultBasket.PutNotStrict(stone, nil)
	} else {
		DefaultBasket.PutNotStrict(stone, value[0])
	}
}

// add a stone to the default basket with a name
func Add(name string, stone Stone, value ...interface{}) {
	if len(value) == 0 {
		DefaultBasket.Add(name, stone)
	} else {
		DefaultBasket.AddWithValue(name, stone, value[0], false)
	}
}

// put a stone into the default basket
func Put(stone Stone, value ...interface{}) {
	if len(value) == 0 {
		DefaultBasket.Put(stone)
	} else {
		DefaultBasket.PutWithValue(stone, value[0], false)
	}
}

// get a stone with the name and the type
func GetStone(name string, t reflect.Type) (stone Stone) {
	return DefaultBasket.GetStone(name, t)
}

// get a tone with the name
func GetStoneWithName(name string) (stone Stone) {
	return DefaultBasket.GetStoneWithName(name)
}

func GetStoneByType(typ interface{}) (stone Stone) {
	return DefaultBasket.GetStoneByType(typ)
}

// register a plugin to basket
func PluginRegister(p Plugin, pt PluginWorkTime) {
	DefaultBasket.PluginRegister(p, pt)
}

// just start
func Start() {
	DefaultBasket.Start()
}

// start and wait interrupt signal to shutdown
func Work() {
	DefaultBasket.Start()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Println("SHUTDOWN:", sig)
			DefaultBasket.ShutDown()
			os.Exit(0)
		}
	}()
}
func Strict() {
	DefaultBasket.Strict()
}
func ShutDown() {
	DefaultBasket.ShutDown()
}

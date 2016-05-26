package summer

import (
	"reflect"
	"strconv"
	"strings"
"qiniupkg.com/x/log.v7"
)

func init() {
	defaultBasket.PluginRegister(&RefPlugin{defaultBasket}, AfterInit)
}

type RefPlugin struct {
	basket Basket
}

func (this *RefPlugin) Look(Holder *Holder, path string) reflect.Value {
	stack := strings.Split(path, ".")
	log.Println("[ref]", path,Holder.class)
	holder := this.basket.NameHolder(stack[0])
	Holder.depends = append(holder.depends, holder)
	root := holder.stone
	value := reflect.ValueOf(root)
	for index, name := range stack {
		if index == 0 {
			continue
		}
		value = this.lookChildren(value, name)
	}
	return value
}

func (this *RefPlugin) lookChildren(parent reflect.Value, childrenName string) reflect.Value {
	pKind := parent.Kind()
	if pKind == reflect.Ptr {
		return this.lookChildren(parent.Elem(), childrenName)
	}
	if pKind == reflect.Struct {
		return parent.FieldByName(childrenName)
	}
	if pKind == reflect.Map {
		return parent.MapIndex(reflect.ValueOf(childrenName))
	}
	if pKind == reflect.Array || pKind == reflect.Slice {
		c, err := strconv.Atoi(childrenName)
		if err != nil {
			panic(err)
		}
		return parent.Index(c)
	}
	panic("sorry i dont know what happend")
}

func (this *RefPlugin) Prefix() string {
	return "$"
}
func (this *RefPlugin) ZIndex() int {
	return 1
}

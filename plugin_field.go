package summer

import (
	"reflect"
	"strconv"
	"strings"
)

type FieldReferencePlugin struct {
}

func (plugin *FieldReferencePlugin) Look(holder *Holder, path string, sf *reflect.StructField) reflect.Value {
	stack := strings.Split(path, ".")
	logger.Debug("[ref]", path, holder.Type)
	foundHolder := holder.Basket.GetStoneHolderWithName(stack[0])
	if foundHolder == nil {
		panic("the " + stack[0] + " not found")
	}
	holder.Dependencies[foundHolder] = true
	root := foundHolder.Stone
	value := reflect.ValueOf(root)
	for index, name := range stack {
		if index == 0 {
			continue
		}
		value = plugin.lookChildren(value, name)
	}
	return value
}

func (plugin *FieldReferencePlugin) lookChildren(parent reflect.Value, childrenName string) reflect.Value {
	pKind := parent.Kind()
	if pKind == reflect.Ptr {
		return plugin.lookChildren(parent.Elem(), childrenName)
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
	panic("sorry i dont know what happended")
}

func (plugin *FieldReferencePlugin) Prefix() string {
	return "$"
}
func (plugin *FieldReferencePlugin) ZIndex() int {
	return 1
}

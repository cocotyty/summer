package features

import (
"reflect"
"strings"
	"strconv"
)

type RefPlugin struct {
	basket Basket
}

func (this *RefPlugin)Look(path string) reflect.Value {
	stack := strings.Split(path, ".")

	root := this.basket.NameStone(stack[0])
	value := reflect.ValueOf(root)
	for index, name := range stack {
		if index == 0 {
			continue
		}
		value = this.lookChildren(value, name)
	}
	return value
}
func (this *RefPlugin)lookChildren(parent reflect.Value, childrenName string) reflect.Value {
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

func (this *RefPlugin)Prefix() string {
	return "$"
}

func (this *RefPlugin)ZIndex() int {
	return 1
}
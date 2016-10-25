package summer

import (
	"errors"
	"reflect"
	"strings"
)

// a holder that can hold stone
type Holder struct {
	Stone           Stone
	Class           reflect.Type
	PointerClass    reflect.Type
	Value           reflect.Value
	Basket          *Basket
	Dependents      []*Holder
	PreTagRootValue interface{}
}

func newHolder(stone Stone, basket *Basket) *Holder {
	if reflect.TypeOf(stone).Kind() == reflect.Func {
		return &Holder{
			Stone:        stone,
			Class:        reflect.TypeOf(stone),
			PointerClass: reflect.TypeOf(stone),
			Value:        reflect.ValueOf(stone),
			Basket:       basket,
			Dependents:   []*Holder{},
		}
	}
	return &Holder{
		Stone:        stone,
		Class:        reflect.TypeOf(stone).Elem(),
		PointerClass: reflect.TypeOf(stone),
		Value:        reflect.ValueOf(stone).Elem(),
		Basket:       basket,
		Dependents:   []*Holder{},
	}
}
func (holder *Holder) ResolveDirectlyDependents() {
	logger.Debug("ResolveDirectlyDependents", holder.Value)
	if holder.Class.Kind() == reflect.Func {
		return
	}
	num := holder.Value.NumField() - 1
	for ; num >= 0; num-- {
		holder.SetDirectDependValue(holder.Value.Field(num), holder.Class.Field(num))
	}
}

// in this step we try to find the stone which the field need
func (holder *Holder) SetDirectDependValue(fieldValue reflect.Value, fieldInfo reflect.StructField) {
	// get the field's tag which belongs to summer
	tag := fieldInfo.Tag.Get("sm")
	if tag == "" {
		return
	}
	if holder.PreTagRootValue != nil {
		tag = preTag(holder.PreTagRootValue, tag)
	}
	logger.Debug("[build Field]", holder.Class.Name(), fieldInfo.Name, fieldInfo.Type.Name(), fieldInfo.Tag, tag)

	// convert text to summer tag option
	tagOption := buildTagOptions(tag)
	// if the field not a straight depend
	if !tagOption.depend {
		// may be the plugin will help it
		holder.Basket.PutDelayField(&DelayField{fieldValue, fieldInfo, tagOption, holder})
		logger.Debug(holder.Class.Name(), " the field [", fieldInfo.Name, "] will be delay. ", tagOption)
		return
	}
	// get stone's name which the field wanted
	var name string
	if tagOption.auto {
		name = fieldInfo.Name
		name = strings.ToLower(name[:1]) + name[1:]
	} else {
		name = tagOption.name
	}
	// get the field type
	fieldType := fieldValue.Type()
	// find the needed stone holder from basket
	hd := holder.Basket.GetStoneHolder(name, fieldType)
	// if holder not found
	if hd == nil {
		// maybe the name is wrong,we suggest the type'name is the stone's name
		if fieldType.Kind() == reflect.Ptr {
			name = fieldType.Elem().Name()
		} else {
			name = fieldType.Name()
		}
		name = strings.ToLower(name[:1]) + name[1:]
		hd = holder.Basket.GetStoneHolder(name, fieldType)
		if hd == nil {
			// we don't know what happened ,maybe you forget put the stone into the basket
			// so just panic
			panic(errors.New("sorry,stone's dependency missed: " + fieldInfo.Name + ",type " + fieldType.Name()))
		}
	}
	// don't forget to record the dependency of the stone we need
	holder.Dependents = append(holder.Dependents, hd)
	fieldValue.Set(reflect.ValueOf(hd.Stone))
	logger.Debug(holder.Class.Name(), " depend on ", hd.Class.Name())
}
func (holder *Holder) init(holders map[*Holder]bool) {
	if stone, ok := holder.Stone.(Init); ok {
		if holders[holder] {
			return
		}
		holders[holder] = true
		for _, v := range holder.Dependents {
			v.init(holders)
		}
		stone.Init()
	}
}
func (holder *Holder) ready(holders map[*Holder]bool) {
	if stone, ok := holder.Stone.(Ready); ok {
		if holders[holder] {
			return
		}
		holders[holder] = true
		for _, v := range holder.Dependents {
			v.ready(holders)
		}
		stone.Ready()
	}
}
func (this *Holder) destroy(holders map[*Holder]bool) {
	if stone, ok := this.Stone.(Destroy); ok {
		if holders[this] {
			return
		}
		holders[this] = true
		for _, v := range this.Dependents {
			v.destroy(holders)
		}
		stone.Destroy()
	}
}

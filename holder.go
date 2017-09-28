package summer

import (
	"errors"
	"reflect"
	"strings"
)

// a holder that can hold stone
type Holder struct {
	IgnoreStrict         bool
	Stone                Stone
	Type                 reflect.Type
	PointerType          reflect.Type
	Value                reflect.Value
	Basket               *Basket
	Dependencies         map[*Holder]bool
	TagTemplateRootValue interface{}
}

func newHolder(stone Stone, basket *Basket) *Holder {
	if reflect.TypeOf(stone).Kind() == reflect.Func {
		return &Holder{
			Stone:        stone,
			Type:         reflect.TypeOf(stone),
			PointerType:  reflect.TypeOf(stone),
			Value:        reflect.ValueOf(stone),
			Basket:       basket,
			Dependencies: map[*Holder]bool{},
		}
	}
	return &Holder{
		Stone:        stone,
		Type:         reflect.TypeOf(stone).Elem(),
		PointerType:  reflect.TypeOf(stone),
		Value:        reflect.ValueOf(stone).Elem(),
		Basket:       basket,
		Dependencies: map[*Holder]bool{},
	}
}
func (h *Holder) ResolveDirectlyDependents() {
	logger.Debug("ResolveDirectlyDependents", h.Value)
	if h.Type.Kind() == reflect.Func {
		return
	}
	num := h.Value.NumField() - 1
	for ; num >= 0; num-- {
		h.SetDirectDependValue(h.Value.Field(num), h.Type.Field(num))
	}
}

// in this step we try to find the stone which the field need
func (h *Holder) SetDirectDependValue(fieldValue reflect.Value, field reflect.StructField) {
	// get the field's tag which belongs to summer
	tag := field.Tag.Get("sm")
	if tag == "" {
		if (!h.IgnoreStrict) && h.Basket.strict && fieldValue.CanSet() {
			panic(" strict mode not support exported field not use summer tag \n" + h.Type.PkgPath() + " " + h.Type.String() + " " + field.Name)
		}
		return
	}
	if tag == "-" {
		return
	}
	if h.TagTemplateRootValue != nil {
		tag = tagTemplateExecute(h.TagTemplateRootValue, tag)
		logger.Debug("[pretag Field]", tag)
	}
	logger.Debug("[build Field]", h.Type.Name(), field.Name, field.Type.Name(), field.Tag, tag)

	// convert text to summer tag option
	tagOption := parseTagOption(tag)
	// if the field not a straight depend
	if !tagOption.depend {
		// may be the plugin will help it
		h.Basket.PutDelayField(&DelayedField{fieldValue, field, tagOption, h})
		logger.Debug(h.Type.Name(), " the field [", field.Name, "] will be delay. ", tagOption)
		return
	}
	// get stone's name which the field wanted
	var name string
	if tagOption.auto {
		name = field.Name
		name = strings.ToLower(name[:1]) + name[1:]
	} else {
		name = tagOption.name
	}
	// get the field type
	fieldType := fieldValue.Type()
	// find the needed stone holder from basket
	hd := h.Basket.GetStoneHolder(name, fieldType)
	// if holder not found
	if hd == nil {
		// maybe the name is wrong,we suggest the type'name is the stone's name
		if fieldType.Kind() == reflect.Ptr {
			name = fieldType.Elem().Name()
		} else {
			name = fieldType.Name()
		}
		name = strings.ToLower(name[:1]) + name[1:]
		hd = h.Basket.GetStoneHolder(name, fieldType)
		if hd == nil {
			// we don't know what happened ,maybe you forget put the stone into the basket
			// so if fieldType's kind  is pointer of struct ,  put a new zero value in basket
			if fieldType.Kind() == reflect.Ptr && fieldType.Elem().Kind() == reflect.Struct {
				h.Basket.Put(reflect.New(fieldType.Elem()).Interface())
				hd = h.Basket.GetStoneHolder(name, fieldType)
				hd.ResolveDirectlyDependents()
			}
			if hd == nil {
				// just panic
				if fieldType.Kind() == reflect.Ptr {
					panic(errors.New("Sorry,stone's dependency missed: " + h.PointerType.String() + " [field] " + field.Name + " [type] pointer of " + fieldType.Elem().PkgPath() + "/" + fieldType.Elem().Name()))
				} else {
					panic(errors.New("Sorry,stone's dependency missed: " + h.PointerType.String() + " [field] " + field.Name + " [type] " + fieldType.PkgPath() + "/" + fieldType.Name()))
				}
			}
		}
	}
	// don't forget to record the dependency of the stone we need
	h.Dependencies[hd] = true
	if fieldValue.CanSet() {
		fieldValue.Set(reflect.ValueOf(hd.Stone))
	} else {
		panic(h.Type.Name() + " depend on " + hd.Type.Name() + ": but not exported value")
	}
	logger.Debug(h.Type.Name(), " depend on ", hd.Type.Name())
}
func (h *Holder) init(holders map[*Holder]bool) {
	if holders[h] {
		return
	}
	holders[h] = true
	for v := range h.Dependencies {
		v.init(holders)
	}
	if stone, ok := h.Stone.(Init); ok {
		stone.Init()
	}
}
func (h *Holder) ready(holders map[*Holder]bool) {
	if holders[h] {
		return
	}
	holders[h] = true
	for v := range h.Dependencies {
		v.ready(holders)
	}
	if stone, ok := h.Stone.(Ready); ok {
		stone.Ready()
	}
}
func (h *Holder) destroy(holders map[*Holder]bool) {
	if holders[h] {
		return
	}
	holders[h] = true
	for v := range h.Dependencies {
		v.destroy(holders)
	}
	if stone, ok := h.Stone.(Destroy); ok {
		stone.Destroy()
	}
}

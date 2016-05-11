package features

import (
	"reflect"
	"strings"
	"qiniupkg.com/x/log.v7"
	"errors"
)

var NotSupportStructErr = errors.New("sorry we not support struct now")
var CannotResolveDependencyErr = errors.New("sorry,stone's dependency missed")

type holder struct {
	stone        Stone
	class        reflect.Type
	pointerClass reflect.Type
	value        reflect.Value
	success      bool
	basket       *basket
}
type tagOption struct {
	auto   bool
	depend bool
	name   string
	path   string
}

func buildTagOptions(tag string) *tagOption {
	to := &tagOption{}
	if tag == "*" {
		to.auto = true
		return
	}
	if len(tag) <= 1 {
		log.Println("bad tag")
		return
	}
	if tag[1] == '$' {
		to.depend = false
		to.path = tag[1:]
		return
	}
	to.depend = true
	to.name = tag
	return
}
func newHolder(stone   Stone, basket *basket) *holder {
	return &holder{stone, reflect.TypeOf(stone).Elem(), reflect.TypeOf(stone), reflect.ValueOf(stone).Elem(), false, basket}
}
func (this *holder)build() {
	num := this.value.NumField()
	num--;
	for ; num >= 0; num-- {
		this.buildFiled(this.value.Field(num), this.class.Field(num))
	}
}

func (this *holder)buildFiled(filedValue reflect.Value, filedInfo reflect.StructField) {
	tag := filedInfo.Tag.Get("sm")
	log.Println(filedInfo, filedInfo.Tag, tag)

	if tag == "" {
		return
	}
	to := buildTagOptions(tag)
	if to.depend {
		t := filedValue.Type()
		name := filedInfo.Name
		name = strings.ToLower(name[:1]) + name[1:]
		hd := this.basket.holder(name, t)
		if hd == nil {
			if t.Kind() == reflect.Ptr {
				name = t.Elem().Name()
			}else {
				name = t.Name()
			}
			name = strings.ToLower(name[:1]) + name[1:]
			hd = this.basket.holder(name, t)
			if hd == nil {
				panic(CannotResolveDependencyErr)
			}
		}
		filedValue.Set(reflect.ValueOf(hd.stone))
	}else {
		this.basket.laterFills = append(this.basket.laterFills, laterFill{filedValue, filedInfo, to})
	}
}

type laterFill struct {
	filedValue reflect.Value
	filedInfo  reflect.StructField
	tagOption  *tagOption
}
type basket struct {
	kv         map[string][]*holder
	laterFills []*laterFill
	//board map[string]*holder
}

func (this *basket)Add(name string, stone Stone) {
	t := reflect.TypeOf(stone)
	if t.Kind() != reflect.Ptr {
		panic(NotSupportStructErr)

	}
	if holders, found := this.kv[name]; found {
		this.kv[name] = append(holders, newHolder(stone, this))
	}else {
		this.kv[name] = []*holder{newHolder(stone, this)}
	}
}

func (this *basket)Put(stone Stone) {
	t := reflect.TypeOf(stone)
	var name string
	if t.Kind() == reflect.Ptr {
		name = t.Elem().Name()
	}else {
		panic(NotSupportStructErr)
	}
	name = strings.ToLower(name[:1]) + name[1:]
	if types, found := this.kv[name]; found {
		this.kv[name] = append(types, newHolder(stone, this))
	}else {
		this.kv[name] = []*holder{newHolder(stone, this)}
	}
}
func (this *basket)build() {
	for _, holders := range this.kv {
		for _, h := range holders {
			h.build()
		}
	}
	for _,lf:=range this.laterFills{

	}
}
func (this *basket)Stone(name string, t reflect.Type) (stone Stone) {
	if holder, found := this.kv[name]; found {
		for _, h := range holder {
			if stone, has := this.find(t, h); has {
				return stone
			}
		}
	}
	for _, holder := range this.kv {
		for _, h := range holder {
			if stone, has := this.find(t, h); has {
				return stone
			}
		}
	}
	return nil
}
func (this *basket)holder(name string, t reflect.Type) (h *holder) {
	if holder, found := this.kv[name]; found {
		for _, h := range holder {
			if _, has := this.find(t, h); has {
				return h
			}
		}
	}
	for _, holder := range this.kv {
		for _, h := range holder {
			if _, has := this.find(t, h); has {
				return h
			}
		}
	}
	return nil
}
func (this *basket)find(t reflect.Type, h *holder) (Stone, bool) {
	if t.Kind() == reflect.Interface {
		if reflect.TypeOf(h.stone).Implements(t) {
			return h.stone, true
		}
		return nil, false
	}
	if t.Kind() == reflect.Struct {
		t = reflect.PtrTo(t)
	}
	if h.pointerClass.AssignableTo(t)&&h.pointerClass.ConvertibleTo(t) {
		return h.stone, true
	}
	return nil, false
}

func (this *basket)Start() {
	this.build()
}
func (this *basket)ShutDown() {

}
type FillPlugin struct {


}
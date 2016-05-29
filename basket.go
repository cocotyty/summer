package summer

import (
	"errors"
	"reflect"
	"sort"
	"strings"
	"qiniupkg.com/x/log.v7"
)

var NotSupportStructErr = errors.New("sorry we not support struct now")
var CannotResolveDependencyErr = errors.New("sorry,stone's dependency missed")

type DelayField struct {
	filedValue reflect.Value
	filedInfo  reflect.StructField
	tagOption  *tagOption
	Holder     *Holder
}

type basket struct {
	kv          map[string][]*Holder
	delayFields map[string][]*DelayField
	plugins     map[PluginWorkTime]pluginList
}

func (this *basket)PutDelayField(field *DelayField) {
	list, has := this.delayFields[field.tagOption.prefix]
	if !has {
		list = []*DelayField{}
	}
	this.delayFields[field.tagOption.prefix] = append(list, field)
}
func NewBasket() Basket {
	return &basket{make(map[string][]*Holder), make(map[string][]*DelayField), make(map[PluginWorkTime]pluginList)}
}
func (this *basket)NameHolder(name string) *Holder {
	if holders, found := this.kv[name]; found {
		return holders[0]
	}
	return nil
}
func (this *basket) Add(name string, stone Stone) {
	t := reflect.TypeOf(stone)
	if t.Kind() != reflect.Ptr {
		panic(NotSupportStructErr)
	}
	if holders, found := this.kv[name]; found {
		this.kv[name] = append(holders, newHolder(stone, this))
	} else {
		this.kv[name] = []*Holder{newHolder(stone, this)}
	}
}
func (this *basket) PluginRegister(plugin Plugin, t PluginWorkTime) {
	log.Println("[plugin register][", plugin.Prefix(), "]", t)
	list, ok := this.plugins[t]
	if !ok {
		list = pluginList{}
	}
	list = append(list, plugin)
	this.plugins[t] = list
}
func (this *basket) Put(stone Stone) {
	t := reflect.TypeOf(stone)
	var name string
	if t.Kind() == reflect.Ptr {
		name = t.Elem().Name()
	} else {
		panic(NotSupportStructErr)
	}
	name = strings.ToLower(name[:1]) + name[1:]
	if types, found := this.kv[name]; found {
		this.kv[name] = append(types, newHolder(stone, this))
	} else {
		this.kv[name] = []*Holder{newHolder(stone, this)}
	}
}
func (this *basket) ResolveStonesDirectlyDependents() {
	this.Each(func(holder *Holder) {
		holder.ResolveDirectlyDependents()
	})
}
func (this *basket) pluginWorks(worktime PluginWorkTime) {
	log.Println("[plugin][start-tag-map]")
	sort.Sort(this.plugins[worktime])
	// choose which plugins will work at this worktime
	list := this.plugins[worktime]
	for _, plugin := range list {
		log.Println("[plugin][load][", worktime, "]:", plugin.Prefix())
		delayList := this.delayFields[plugin.Prefix()]
		for _, field := range delayList {
			this.pluginWork(plugin, field)
		}
	}
	log.Println("[plugin][finish]")
}
func (this *basket)pluginWork(plugin Plugin, field *DelayField) {
	// find the value we need from plugin
	foundValue := plugin.Look(field.Holder, field.tagOption.path)
	// verify value
	if !foundValue.IsValid() {
		log.Error(plugin.Prefix(), ".", field.tagOption.path, " not found")
		return
	}
	// verify if the field can set a value
	if !field.filedValue.CanSet() {
		log.Error("can not set the value ", field.filedInfo.Name, " tag:", field.filedInfo.Tag, ",may be an unexported value ")
		return
	}
	log.Println("[plugin][path]", field.Holder.Class, field.tagOption.path, foundValue.Interface())
	if field.filedInfo.Type.Kind() == foundValue.Kind() {
		field.filedValue.Set(foundValue)
		return
	}
	if field.filedInfo.Type.Kind() == reflect.Ptr && foundValue.Kind() != reflect.Ptr {
		field.filedValue.Set(foundValue.Addr())
		return
	}
	if field.filedInfo.Type.Kind() != reflect.Ptr && foundValue.Kind() == reflect.Ptr {
		field.filedValue.Set(foundValue.Elem())
		return
	}
	if ( field.filedInfo.Type.Kind() == reflect.Int || field.filedInfo.Type.Kind() == reflect.Int8 || field.filedInfo.Type.Kind() == reflect.Int16 || field.filedInfo.Type.Kind() == reflect.Int32 || field.filedInfo.Type.Kind() == reflect.Int64) {
		switch value := foundValue.Interface().(type) {
		case int:
			field.filedValue.SetInt(int64(value))
		case int8:
			field.filedValue.SetInt(int64(value))
		case int16:
			field.filedValue.SetInt(int64(value))
		case int32:
			field.filedValue.SetInt(int64(value))
		case int64:
			field.filedValue.SetInt(int64(value))
		default:
			log.Error("can not set the value ", field.filedInfo.Name, " tag:", field.filedInfo.Tag, " because ", field.filedInfo.Type, "!=", foundValue.Kind())
		}
		return
	}
	if ( foundValue.Kind() == reflect.Float64 || foundValue.Kind() == reflect.Float32 ) {
		switch value := foundValue.Interface().(type) {
		case float32:
			field.filedValue.SetFloat(float64(value))
		case float64:
			field.filedValue.SetFloat(float64(value))
		default:
			log.Error("can not set the value ", field.filedInfo.Name, " tag:", field.filedInfo.Tag, " because ", field.filedInfo.Type, "!=", foundValue.Kind())
		}
		return
	}
	log.Error("can not set the value ", field.filedInfo.Name, " tag:", field.filedInfo.Tag, " because ", field.filedInfo.Type, "!=", foundValue.Kind())
}
func (this *basket) Stone(name string, t reflect.Type) (stone Stone) {
	if holder, found := this.kv[name]; found {
		for _, h := range holder {
			if stone, has := this.findStone(t, h); has {
				return stone
			}
		}
	}
	for _, holder := range this.kv {
		for _, h := range holder {
			if stone, has := this.findStone(t, h); has {
				return stone
			}
		}
	}
	return nil
}
func (this *basket) NameStone(name string) (stone Stone) {
	if holders, found := this.kv[name]; found {
		return holders[0].Stone
	}
	return nil
}
func (this *basket) findHolder(name string, t reflect.Type) (h *Holder) {
	if holder, found := this.kv[name]; found {
		for _, h := range holder {
			if _, has := this.findStone(t, h); has {
				return h
			}
		}
	}
	for _, holder := range this.kv {
		for _, h := range holder {
			if _, has := this.findStone(t, h); has {
				return h
			}
		}
	}
	return nil
}
func (this *basket) findStone(t reflect.Type, h *Holder) (Stone, bool) {
	if t.Kind() == reflect.Interface {
		if reflect.TypeOf(h.Stone).Implements(t) {
			return h.Stone, true
		}
		return nil, false
	}
	if t.Kind() == reflect.Struct {
		t = reflect.PtrTo(t)
	}
	if h.PointerClass.AssignableTo(t) && h.PointerClass.ConvertibleTo(t) {
		return h.Stone, true
	}
	return nil, false
}
func (this *basket) Start() {
	this.ResolveStonesDirectlyDependents()
	this.pluginWorks(BeforeInit)
	this.initStones()
	this.pluginWorks(AfterInit)
	this.tellStoneReady()
}
func (this *basket)initStones() {
	set := map[*Holder]bool{}
	this.Each(func(holder *Holder) {
		holder.init(set)
	})
}
func (this *basket)tellStoneReady() {
	set := map[*Holder]bool{}
	this.Each(func(holder *Holder) {
		holder.ready(set)
	})
}
func (this *basket) ShutDown() {
	set := map[*Holder]bool{}
	this.Each(func(holder *Holder) {
		holder.destroy(set)
	})
}
func (this *basket)Each(fn func(holder *Holder)) {
	for _, holders := range this.kv {
		for _, holder := range holders {
			fn(holder)
		}
	}
}
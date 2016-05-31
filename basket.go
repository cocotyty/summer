package summer

import (
	"errors"
	"reflect"
	"sort"
	"strings"
	"qiniupkg.com/x/log.v7"
)

var NotSupportStructErr = errors.New("sorry we not support struct now")
var NotSupportContainsDot = errors.New("sorry we not support name contains a dot")
var CannotResolveDependencyErr = errors.New("sorry,stone's dependency missed")

type Basket struct {
	kv          map[string][]*Holder
	delayFields map[string][]*DelayField
	plugins     map[PluginWorkTime]pluginList
}

func (this *Basket)PutDelayField(field *DelayField) {
	list, has := this.delayFields[field.tagOption.prefix]
	if !has {
		list = []*DelayField{}
	}
	this.delayFields[field.tagOption.prefix] = append(list, field)
}
func NewBasket() *Basket {
	return &Basket{make(map[string][]*Holder), make(map[string][]*DelayField), make(map[PluginWorkTime]pluginList)}
}
// add a stone to basket,the stone must be struct's pointer
func (this *Basket) Add(name string, stone Stone) {
	if strings.Contains(name, ".") {
		panic(NotSupportContainsDot)
	}
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
// put a stone into basket ,the stone must be struct's pointer,the stone name will be that's type's name with first character lowercase
// for example,if stone's type is Foo then the stone will get a name that is "foo"
func (this *Basket) Put(stone Stone) {
	t := reflect.TypeOf(stone)
	var name string
	if t.Kind() == reflect.Ptr {
		name = t.Elem().Name()
	} else {
		panic(NotSupportStructErr)
	}
	name = strings.ToLower(name[:1]) + name[1:]
	log.Debug("regitor ", name)
	if types, found := this.kv[name]; found {
		this.kv[name] = append(types, newHolder(stone, this))
	} else {
		this.kv[name] = []*Holder{newHolder(stone, this)}
	}
}
// register a plugin to basket
func (this *Basket) PluginRegister(plugin Plugin, t PluginWorkTime) {
	log.Debug("[plugin register][", plugin.Prefix(), "]", t)
	list, ok := this.plugins[t]
	if !ok {
		list = pluginList{}
	}
	list = append(list, plugin)
	this.plugins[t] = list
}

func (this *Basket) resolveStonesDirectlyDependents() {
	this.Each(func(holder *Holder) {
		holder.ResolveDirectlyDependents()
	})
}
func (this *Basket) pluginWorks(worktime PluginWorkTime) {
	log.Debug("[plugin][start-tag-map]")
	sort.Sort(this.plugins[worktime])
	// choose which plugins will work at this worktime
	list := this.plugins[worktime]
	for _, plugin := range list {
		log.Debug("[plugin][load][", worktime, "]:", plugin.Prefix())
		delayList := this.delayFields[plugin.Prefix()]
		for _, field := range delayList {
			this.pluginWork(plugin, field)
		}
	}
	log.Debug("[plugin][finish]")
}
func (this *Basket) pluginWork(plugin Plugin, field *DelayField) {
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
	log.Debug("[plugin][path]", field.Holder.Class, field.tagOption.path, foundValue.Interface())
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
// get a stone from basket
func (this *Basket) GetStone(name string, t reflect.Type) (stone Stone) {
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
// get a stone from basket
func (this *Basket) GetStoneWithName(name string) (stone Stone) {
	if holders, found := this.kv[name]; found {
		return holders[0].Stone
	}
	return nil
}
// get a stone holder from basket
func (this *Basket) GetStoneHolder(name string, t reflect.Type) (h *Holder) {
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
// get a stone holder from basket
func (this *Basket) GetStoneHolderWithName(name string) *Holder {
	if holders, found := this.kv[name]; found {
		return holders[0]
	}
	return nil
}
func (this *Basket) findStone(t reflect.Type, h *Holder) (Stone, bool) {
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
// start work
//
// # 1 : summer will resolve the direct dependency
//
// # 2 : summer will call all stones's Init method<br/>
// if a directly depend on b, b directly depend on c and d ,then a will init after b init,and b will after c and d
//
// # 3 : summer will call all stones's Ready method<br/>
// if a  depend on b, b  depend on c and d ,then a will init after b init,and b will after c and d
//
func (this *Basket) Start() {
	this.resolveStonesDirectlyDependents()
	this.pluginWorks(BeforeInit)
	this.initStones()
	this.pluginWorks(AfterInit)
	this.tellStoneReady()
}
func (this *Basket)initStones() {
	set := map[*Holder]bool{}
	this.Each(func(holder *Holder) {
		holder.init(set)
	})
}
func (this *Basket)tellStoneReady() {
	set := map[*Holder]bool{}
	this.Each(func(holder *Holder) {
		holder.ready(set)
	})
}
// shutdown will call all stone's Destroy method
func (this *Basket) ShutDown() {
	set := map[*Holder]bool{}
	this.Each(func(holder *Holder) {
		holder.destroy(set)
	})
}
func (this *Basket)Each(fn func(holder *Holder)) {
	for _, holders := range this.kv {
		for _, holder := range holders {
			fn(holder)
		}
	}
}
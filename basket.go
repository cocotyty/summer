package summer

import (
	"errors"
	"reflect"
	"sort"
	"strings"
)

var NotSupportStructErr = errors.New("sorry we not support struct now")
var NotSupportContainsDot = errors.New("sorry we not support name contains a dot")

type Basket struct {
	// strict mode
	strict         bool
	nameHoldersMap map[string][]*Holder
	delayedFields  map[string][]*DelayedField
	plugins        map[PluginWorkTime]pluginList
}

func (basket *Basket) Strict() {
	basket.strict = true
}
func (basket *Basket) PutDelayField(field *DelayedField) {
	list, has := basket.delayedFields[field.tagOption.prefix]
	if !has {
		list = []*DelayedField{}
	}
	basket.delayedFields[field.tagOption.prefix] = append(list, field)
}
func NewBasket() *Basket {
	return &Basket{
		strict:         false,
		nameHoldersMap: make(map[string][]*Holder),
		delayedFields:  make(map[string][]*DelayedField),
		plugins:        make(map[PluginWorkTime]pluginList)}
}
func (basket *Basket) AddNotStrict(name string, stone Stone, value interface{}) {
	basket.AddWithValue(name, stone, value, true)
}
func (basket *Basket) PutNotStrict(stone Stone, value interface{}) {
	basket.PutWithValue(stone, value, true)
}

// add a stone to basket,the stone must be struct's pointer
func (basket *Basket) Add(name string, stone Stone) {
	basket.AddWithValue(name, stone, nil, false)
}
func (basket *Basket) AddWithValue(name string, stone Stone, root interface{}, ignoreStrict bool) {
	if strings.Contains(name, ".") {
		panic(NotSupportContainsDot)
	}
	t := reflect.TypeOf(stone)
	storeKind := t.Kind()
	if storeKind != reflect.Ptr && storeKind != reflect.Func {
		panic(NotSupportStructErr)
	}
	holder := newHolder(stone, basket)
	holder.IgnoreStrict = ignoreStrict
	if holders, found := basket.nameHoldersMap[name]; found {
		basket.nameHoldersMap[name] = append(holders, holder)
	} else {
		basket.nameHoldersMap[name] = []*Holder{holder}
	}
	holder.TagTemplateRootValue = root
}

func (basket *Basket) PutWithValue(stone Stone, root interface{}, ignoreStrict bool) {
	t := reflect.TypeOf(stone)
	var name string
	if t.Kind() == reflect.Ptr {
		name = t.Elem().Name()
	} else if t.Kind() == reflect.Func {
		name = "funcs"
	} else {
		panic(NotSupportStructErr)
	}
	name = strings.ToLower(name[:1]) + name[1:]
	logger.Debug("registor ", name)
	basket.AddWithValue(name, stone, root, ignoreStrict)
}

// put a stone into basket ,the stone must be struct's pointer,the stone name will be that's type's name with first character lowercase
// for example,if stone's type is Foo then the stone will get a name that is "foo"
func (basket *Basket) Put(stone Stone) {
	basket.PutWithValue(stone, nil, false)
}

// register a plugin to basket
func (basket *Basket) PluginRegister(plugin Plugin, t PluginWorkTime) {
	logger.Debug("[plugin register][", plugin.Prefix(), "]", t)
	list, ok := basket.plugins[t]
	if !ok {
		list = pluginList{}
	}
	list = append(list, plugin)
	basket.plugins[t] = list
}

func (basket *Basket) resolveStonesDirectlyDependents() {
	basket.SafeEach(func(holder *Holder) {
		holder.ResolveDirectlyDependents()
	})
}
func (basket *Basket) pluginWorks(wt PluginWorkTime) {
	logger.Debug("[plugin][start-tag-map]")
	sort.Sort(basket.plugins[wt])
	// choose which plugins will work at this worktime
	list := basket.plugins[wt]
	for _, plugin := range list {
		logger.Debug("[plugin][load][", wt, "]:", plugin.Prefix())
		delayList := basket.delayedFields[plugin.Prefix()]
		for _, field := range delayList {
			basket.pluginWork(plugin, field)
		}
	}
	logger.Debug("[plugin][finish]")
}
func (basket *Basket) pluginWork(plugin Plugin, field *DelayedField) {
	// find the value we need from plugin
	foundValue := plugin.Look(field.holder, field.tagOption.path, &field.field)
	// verify value
	if !foundValue.IsValid() {
		logger.Error(plugin.Prefix(), ".", field.tagOption.path, " not found")
		return
	}
	// verify if the field can set a value
	if !field.value.CanSet() {
		logger.Error("can not set the value ", field.field.Name, " tag:", field.field.Tag, ",may be an unexported value ")
		return
	}
	logger.Debug("[plugin][path]", field.holder.Type, field.tagOption.path, foundValue.Interface())
	if field.field.Type.Kind() == foundValue.Kind() {
		field.value.Set(foundValue)
		return
	}
	if field.field.Type.Kind() == reflect.Interface {
		if foundValue.Type().AssignableTo(field.field.Type) && foundValue.Type().ConvertibleTo(field.field.Type) {
			field.value.Set(foundValue)
			return
		}
	}
	if field.field.Type.Kind() == reflect.Ptr && foundValue.Kind() != reflect.Ptr {
		field.value.Set(foundValue.Addr())
		return
	}
	if field.field.Type.Kind() != reflect.Ptr && foundValue.Kind() == reflect.Ptr {
		field.value.Set(foundValue.Elem())
		return
	}
	if field.field.Type.Kind() == reflect.Int || field.field.Type.Kind() == reflect.Int8 || field.field.Type.Kind() == reflect.Int16 || field.field.Type.Kind() == reflect.Int32 || field.field.Type.Kind() == reflect.Int64 {
		switch value := foundValue.Interface().(type) {
		case int:
			field.value.SetInt(int64(value))
		case int8:
			field.value.SetInt(int64(value))
		case int16:
			field.value.SetInt(int64(value))
		case int32:
			field.value.SetInt(int64(value))
		case int64:
			field.value.SetInt(int64(value))
		default:
			logger.Error("can not set the value ", field.field.Name, " tag:", field.field.Tag, " because ", field.field.Type, "!=", foundValue.Kind())
		}
		return
	}
	if foundValue.Kind() == reflect.Float64 || foundValue.Kind() == reflect.Float32 {
		switch value := foundValue.Interface().(type) {
		case float32:
			field.value.SetFloat(float64(value))
		case float64:
			field.value.SetFloat(float64(value))
		default:
			logger.Error("can not set the value ", field.field.Name, " tag:", field.field.Tag, " because ", field.field.Type, "!=", foundValue.Kind())
		}
		return
	}
	logger.Error("can not set the value ", field.field.Name, " tag:", field.field.Tag, " because ", field.field.Type, "!=", foundValue.Kind())
}

// get a stone from basket
func (basket *Basket) GetStone(name string, t reflect.Type) (stone Stone) {
	if holder, found := basket.nameHoldersMap[name]; found {
		for _, h := range holder {
			if stone, has := basket.findStone(t, h); has {
				return stone
			}
		}
	}
	for _, holder := range basket.nameHoldersMap {
		for _, h := range holder {
			if stone, has := basket.findStone(t, h); has {
				return stone
			}
		}
	}
	return nil
}

// get a stone from basket
func (basket *Basket) GetStoneWithName(name string) (stone Stone) {
	if holders, found := basket.nameHoldersMap[name]; found {
		return holders[0].Stone
	}
	return nil
}
func (basket *Basket) GetStoneByType(typ interface{}) (result Stone) {
	t := reflect.TypeOf(typ)
	if t.Kind() != reflect.Ptr && t.Kind() != reflect.Func {
		t = reflect.New(t).Type()
	}

	basket.EachHolder(func(name string, holder *Holder) bool {
		if holder.PointerType == t {
			result = holder.Stone
			return true
		}
		return false
	})
	return
}

// get a stone holder from basket
func (basket *Basket) GetStoneHolder(name string, t reflect.Type) (h *Holder) {
	if holder, found := basket.nameHoldersMap[name]; found {
		for _, h := range holder {
			if _, has := basket.findStone(t, h); has {
				return h
			}
		}
	}
	for _, holder := range basket.nameHoldersMap {
		for _, h := range holder {
			if _, has := basket.findStone(t, h); has {
				return h
			}
		}
	}
	return nil
}

// get a stone holder from basket
func (basket *Basket) GetStoneHolderWithName(name string) *Holder {
	if holders, found := basket.nameHoldersMap[name]; found {
		return holders[0]
	}
	return nil
}
func (basket *Basket) findStone(t reflect.Type, h *Holder) (Stone, bool) {
	logger.Debug(t.PkgPath(), t, h.Type)
	if t.Kind() == reflect.Interface {
		if reflect.TypeOf(h.Stone).Implements(t) {
			return h.Stone, true
		}
		return nil, false
	}
	if t.Kind() == reflect.Struct {
		t = reflect.PtrTo(t)
	}
	if h.PointerType.AssignableTo(t) && h.PointerType.ConvertibleTo(t) {
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
func (basket *Basket) Start() {
	basket.resolveStonesDirectlyDependents()
	basket.pluginWorks(BeforeInit)
	basket.initStones()
	basket.pluginWorks(AfterInit)
	basket.tellStoneReady()
}
func (basket *Basket) initStones() {
	set := map[*Holder]bool{}
	basket.Each(func(holder *Holder) {
		holder.init(set)
	})
}
func (basket *Basket) tellStoneReady() {
	set := map[*Holder]bool{}
	basket.Each(func(holder *Holder) {
		holder.ready(set)
	})
}

// shutdown will call all stone's Destroy method
func (basket *Basket) ShutDown() {
	set := map[*Holder]bool{}
	basket.Each(func(holder *Holder) {
		holder.destroy(set)
	})
}
func (basket *Basket) copiedNameHoldersMap() map[string][]*Holder {
	m := map[string][]*Holder{}
	for name, holders := range basket.nameHoldersMap {
		copyHolders := make([]*Holder, len(holders))
		copy(copyHolders, holders)
		m[name] = copyHolders
	}
	return m
}
func (basket *Basket) SafeEach(fn func(holder *Holder)) {
	m := basket.copiedNameHoldersMap()
	for _, holders := range m {
		for _, holder := range holders {
			fn(holder)
		}
	}
}
func (basket *Basket) Each(fn func(holder *Holder)) {
	for _, holders := range basket.nameHoldersMap {
		for _, holder := range holders {
			fn(holder)
		}
	}
}
func (basket *Basket) EachHolder(fn func(name string, holder *Holder) bool) {
	for name, holders := range basket.nameHoldersMap {
		for _, holder := range holders {
			if fn(name, holder) {
				return
			}
		}
	}
}

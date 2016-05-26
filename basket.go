package summer

import (
	"errors"
	"reflect"
	"sort"
	"strings"
	"qiniupkg.com/x/log.v7"
	"github.com/pelletier/go-toml"
	"fmt"
)

func TomlFile(path string) error {
	tree, err := toml.LoadFile(path)
	if err != nil {
		return err
	}
	defaultBasket.PluginRegister(&TomlPlugin{tree}, BeforeInit)
	return nil
}

func Toml(src string) error {
	tree, err := toml.Load(src)
	if err != nil {
		return err
	}
	defaultBasket.PluginRegister(&TomlPlugin{tree}, BeforeInit)
	return nil
}

var defaultBasket = NewBasket()

func Add(name string, stone Stone) {
	defaultBasket.Add(name, stone)
}
func Put(stone Stone) {
	defaultBasket.Put(stone)
}
func GetStone(name string, t reflect.Type) (stone Stone) {
	return defaultBasket.Stone(name, t)

}
func NameStone(name string) (stone Stone) {
	return defaultBasket.NameStone(name)
}
func PluginRegister(p Plugin, pt PluginWorkTime) {
	defaultBasket.PluginRegister(p, pt)
}
func Start() {
	defaultBasket.Start()
}
func ShutDown() {
	defaultBasket.ShutDown()
}

var NotSupportStructErr = errors.New("sorry we not support struct now")
var CannotResolveDependencyErr = errors.New("sorry,stone's dependency missed")

type Holder struct {
	stone        Stone
	class        reflect.Type
	pointerClass reflect.Type
	value        reflect.Value
	success      bool
	basket       *basket
	depends      []*Holder
}
type tagOption struct {
	auto   bool
	depend bool
	name   string
	path   string
	prefix string
}

func buildTagOptions(tag string) *tagOption {
	to := &tagOption{}
	if tag == "*" {
		to.depend = true
		to.auto = true
		return to
	}
	if len(tag) <= 1 {
		log.Error("bad tag :", tag)
		return to
	}
	if strings.Contains(tag, ".") {
		to.depend = false
		to.prefix = tag[:strings.Index(tag, ".")]
		to.path = tag[strings.Index(tag, ".") + 1:]
		return to
	}
	to.depend = true
	to.name = tag
	return to
}
func newHolder(stone Stone, basket *basket) *Holder {
	return &Holder{stone, reflect.TypeOf(stone).Elem(), reflect.TypeOf(stone), reflect.ValueOf(stone).Elem(), false, basket, []*Holder{}}
}
func (this *Holder) build() {
	num := this.value.NumField()
	num--
	for ; num >= 0; num-- {
		this.buildFiled(this.value.Field(num), this.class.Field(num))
	}
}

func (this *Holder) buildFiled(filedValue reflect.Value, filedInfo reflect.StructField) {
	tag := filedInfo.Tag.Get("sm")
	fmt.Println("[build filed]", filedInfo.Name, ",", filedInfo.Tag)
	log.Debug("[build filed]", filedInfo, filedInfo.Tag, tag)
	if tag == "" {
		return
	}
	to := buildTagOptions(tag)
	var name string
	if to.auto {
		name = filedInfo.Name
		name = strings.ToLower(name[:1]) + name[1:]
	}else {
		name = to.name
	}
	if to.depend {
		t := filedValue.Type()
		hd := this.basket.holder(name, t)
		if hd == nil {
			if t.Kind() == reflect.Ptr {
				name = t.Elem().Name()
			} else {
				name = t.Name()
			}
			name = strings.ToLower(name[:1]) + name[1:]
			hd = this.basket.holder(name, t)
			if hd == nil {
				panic(CannotResolveDependencyErr)
			}
		}
		filedValue.Set(reflect.ValueOf(hd.stone))
	} else {
		this.basket.laterFills = append(this.basket.laterFills, &laterFill{filedValue, filedInfo, to, this})
	}
}

type laterFill struct {
	filedValue reflect.Value
	filedInfo  reflect.StructField
	tagOption  *tagOption
	Holder     *Holder
}
type pluginList []Plugin

func (list pluginList) Len() int {
	return len(list)
}
func (list pluginList) Less(i, j int) bool {
	return list[i].ZIndex() < list[j].ZIndex()
}
func (list pluginList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

type basket struct {
	kv         map[string][]*Holder
	laterFills []*laterFill
	//board map[string]*holder
	plugins    map[PluginWorkTime]pluginList
}

func NewBasket() Basket {
	return &basket{make(map[string][]*Holder), []*laterFill{}, make(map[PluginWorkTime]pluginList)}
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
	fmt.Println("[plugin register][", plugin.Prefix(), "]", t)
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
func (this *basket) build() {
	for _, holders := range this.kv {
		for _, h := range holders {
			h.build()
		}
	}

}
func (this *basket) pluginWorks(worktime PluginWorkTime) {
	log.Debug("[plugin][start-tag-map]")
	m := map[string][]*laterFill{}
	for _, lf := range this.laterFills {
		if list, has := m[lf.tagOption.prefix]; has {
			list = append(list, lf)
		} else {
			m[lf.tagOption.prefix] = []*laterFill{lf}
		}
	}
	sort.Sort(this.plugins[worktime])
	log.Info("[plugins]", this.plugins[worktime]);
	list := this.plugins[worktime]
	for _, p := range list {
		log.Debug("[plugin][load][", worktime, "]:", p.Prefix())
		laters := m[p.Prefix()]
		for _, l := range laters {
			v := p.Look(l.Holder, l.tagOption.path)
			log.Debug("[plugin][path]", l.tagOption.path, v.Interface())
			if !l.filedValue.CanSet() {
				log.Error("can not set the value ", l.filedInfo.Name, " tag:", l.filedInfo.Tag, ",may be an unexported value ")
				continue
			}
			if l.filedInfo.Type.Kind() != v.Kind() {
				if l.filedInfo.Type.Kind() == reflect.Ptr && v.Kind() != reflect.Ptr {
					l.filedValue.Set(v.Addr())
				}else if l.filedInfo.Type.Kind() != reflect.Ptr && v.Kind() == reflect.Ptr {
					l.filedValue.Set(v.Elem())
				}else if ( l.filedInfo.Type.Kind() == reflect.Int || l.filedInfo.Type.Kind() == reflect.Int8 || l.filedInfo.Type.Kind() == reflect.Int16 || l.filedInfo.Type.Kind() == reflect.Int32 || l.filedInfo.Type.Kind() == reflect.Int64) {
					switch value := v.Interface().(type) {
					case int8:
						l.filedValue.SetInt(int64(value))
					case int16:
						l.filedValue.SetInt(int64(value))
					case int32:
						l.filedValue.SetInt(int64(value))
					case int64:
						l.filedValue.SetInt(int64(value))
					case int:
						l.filedValue.SetInt(int64(value))
					default:
						log.Error("can not set the value ", l.filedInfo.Name, " tag:", l.filedInfo.Tag, " because ", l.filedInfo.Type, "!=", v.Kind())
					}
				}else if ( v.Kind() == reflect.Float64 || v.Kind() == reflect.Float32 ) {
					switch value := v.Interface().(type) {
					case float32:
						l.filedValue.SetFloat(float64(value))
					case float64:
						l.filedValue.SetFloat(float64(value))
					default:
						log.Error("can not set the value ", l.filedInfo.Name, " tag:", l.filedInfo.Tag, " because ", l.filedInfo.Type, "!=", v.Kind())
					}
				}else {
					log.Error("can not set the value ", l.filedInfo.Name, " tag:", l.filedInfo.Tag, " because ", l.filedInfo.Type, "!=", v.Kind())
				}
			} else {
				l.filedValue.Set(v)
			}
		}
	}
	log.Debug("[plugin][finish]")
}
func (this *basket) Stone(name string, t reflect.Type) (stone Stone) {
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
func (this *basket) NameStone(name string) (stone Stone) {
	if holders, found := this.kv[name]; found {
		return holders[0].stone
	}
	return nil
}
func (this *basket) holder(name string, t reflect.Type) (h *Holder) {
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
func (this *basket) find(t reflect.Type, h *Holder) (Stone, bool) {
	if t.Kind() == reflect.Interface {
		if reflect.TypeOf(h.stone).Implements(t) {
			return h.stone, true
		}
		return nil, false
	}
	if t.Kind() == reflect.Struct {
		t = reflect.PtrTo(t)
	}
	if h.pointerClass.AssignableTo(t) && h.pointerClass.ConvertibleTo(t) {
		return h.stone, true
	}
	return nil, false
}

func (this *basket) Start() {
	this.build()
	this.pluginWorks(BeforeInit)
	for _, holders := range this.kv {
		for _, holder := range holders {
			if initer, ok := holder.stone.(Init); ok {
				log.Debug("[init]", holder.class.Name(), holder.stone)
				initer.Init()
			}else {
				log.Debug("[without init]", holder.class.Name(), holder.stone)
			}
		}
	}
	this.pluginWorks(AfterInit)
	set := map[*Holder]bool{}
	for _, holders := range this.kv {
		for _, holder := range holders {
			holder.ready(set)
		}
	}
}
func (this *Holder)ready(holders map[*Holder]bool) {
	if initer, ok := this.stone.(Ready); ok {
		if holders[this] {
			return
		}
		holders[this] = true
		for _, v := range this.depends {
			v.ready(holders)
		}
		initer.Ready()
	}
}
func (this *basket) ShutDown() {
	for _, holders := range this.kv {
		for _, holder := range holders {
			if initer, ok := holder.stone.(Destroy); ok {
				initer.Destroy()
			}
		}
	}
}

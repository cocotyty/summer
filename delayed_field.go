package summer

import "reflect"

type DelayedField struct {
	value     reflect.Value
	field     reflect.StructField
	tagOption *tagOption
	holder    *Holder
}

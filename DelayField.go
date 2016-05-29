package summer

import "reflect"

type DelayField struct {
	filedValue reflect.Value
	filedInfo  reflect.StructField
	tagOption  *tagOption
	Holder     *Holder
}
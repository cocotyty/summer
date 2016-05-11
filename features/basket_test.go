package features

import (
	"testing"
	"github.com/cocotyty/summer/simples"
	"reflect"
)

func TestBasket(t *testing.T) {
	b := &basket{make(map[string][]*holder)}
	b.Put(&simples.A{})
	b.Put(&simples.C{})
	b.Put(&simples.D{})
	b.build()
	b.Stone("A", reflect.TypeOf(&simples.A{})).(*simples.A).Print()
}

package summer

import (
	"fmt"
	"testing"
)

type CycleTypeA struct {
	B *CycleTypeB `sm:"*"`
}
type CycleTypeB struct {
	A   *CycleTypeA `sm:"*"`
	Str string      `sm:"#.str"`
}

type CycleTypeC struct {
	B *CycleTypeB `sm:"*"`
}

func TestStart_Cycle(t *testing.T) {
	a := &CycleTypeA{}
	Put(a)
	Toml(`str="abc"`)
	Start()
	if a.B.Str != "abc" {
		t.Fatal(`a.B.Str != "abc"`)
	}
	if a.B.A != a {
		t.Fatal(`a.B.A != a `)
	}
}

type CommonA struct {
	B *CommonB `sm:"*"`
	C *CommonC `sm:"*"`
}

func (c *CommonA) Init() {
	fmt.Println("CommonA")
}

type CommonB struct {
	Str string
	C   *CommonC `sm:"*"`
	D   *CommonD `sm:"*"`
}

func (c *CommonB) Init() {
	fmt.Println("CommonB")
}

type CommonC struct {
	Str string   `sm:"#.str"`
	D   *CommonD `sm:"*"`
}

func (c *CommonC) Init() {
	fmt.Println("CommonC")
}

type CommonD struct {
}

func (c *CommonD) Init() {
	fmt.Println("CommonD")
}
func TestStart_Common(t *testing.T) {
	a := &CommonA{}
	b := &CommonB{}
	Put(a)
	Put(b)
	Toml(`str="string"`)
	Start()
	if a.B != b {
		t.Fatal(`a.B!=b`)
	}
	if a.C != b.C {
		t.Fatal(`a.C!= b.C`)
	}
	if a.B.D != a.C.D {
		t.Fatal(`a.B.D!= a.C.D`)
	}
	if a.C.Str != "string" {
		t.Fatal(a.C.Str)
	}
}

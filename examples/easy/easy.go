package main

import (
	"github.com/cocotyty/summer"
	"fmt"
)

func init() {
	summer.Toml(`
	[printer]
	prefix="[PRINT]"`)
	summer.Put(&A{})
	summer.Put(&B{})
	summer.Put(&Cat{})
	summer.Put(&Printer{})
	summer.Start()
}

func main() {
	a := summer.GetStoneWithName("a").(*A)
	a.Call()
}

type A struct {
	// $ means you want to get a stone's field , it happened usually after stones inited
	BoyName string `sm:"$.b.Name"`
	// yes,we support interface ,tag is stone's name
	C       C `sm:"cat"`
}

func (a *A)Call() {
	a.C.Print()
	fmt.Println("hi ,I am A")
}

type B struct {
	Name string
}

func (this *B)Init() {
	this.Name = "Boy!"
}

type C interface {
	Print()
}
type Printer struct {
	// if you already set the toml plugin config, you can use the #  ,to get value from toml,
	// # is toml plugin's name
	// toml plugin will work after directly dependency resolved,before init
	Prefix string `sm:"#.printer.prefix"`
}

func (printer *Printer)Print(str string) {
	fmt.Println(printer.Prefix + str)
}

type Cat struct {
	// * is mostly used tag,summer will find by the field's name  or the field's type or both
	Printer *Printer `sm:"*"`
}
func (c *Cat)Ready(){
	fmt.Println("my name is cat,i am ready.")
}
func (c *Cat)Print() {
	c.Printer.Print("Little Cat")
}
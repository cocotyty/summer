# summer
summer like spring ^_&amp;

## What is summer ?

![image](https://github.com/cocotyty/summer/raw/master/logo.png)

Summer is a easy tool if you just want a simple IOC/DI framework. Summer work based on tag "sm".

If you want to use summer,you must put stone to the summer basket first,

and you must call "Start",then summer can work for you.

## What is Stone ?

Stone is go stone,like bean is javabean.

## What is basket ?

A container to store stone,and resolve stones dependents.


## simple use:

```go
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
	summer.Add("lay", &B{})
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
	BoyName string `sm:"$.lay.Name"`
	B       *B `sm:"lay"`
	// yes,we support interface ,tag is stone's name
	C       C `sm:"cat"`
}

func (a *A)Call() {
	a.C.Print()
	fmt.Println("hi ,I am A", "bodys name:", a.BoyName)
	fmt.Println(a.B)
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

func (c *Cat)Ready() {
	fmt.Println("my name is cat,i am ready.")
}
func (c *Cat)Print() {
	c.Printer.Print("Little Cat")
}
```

output:
```text
my name is cat,i am ready.
[PRINT]Little Cat
hi ,I am A bodys name: Boy!
&{Boy!}
```

really * easy *  as you see.

[doc](http://godoc.org/github.com/cocotyty/summer)

if you have questions ,you can send me a email  , my email address  is cocotyty@sina.com 



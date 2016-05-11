package simples

import "qiniupkg.com/x/log.v7"

type A struct {
	Name string
	D    *D `sm:"auto"`
}

func (a *A)Print() {
	log.Println(a.Name, "[d]", a.D)
}

type C struct {
	Name string
	A    *A `sm:"auto"`
}

func (c *C)Print() {
	log.Println(c.Name, "[a]", c.A)
}

type D struct {
	Name string
	C    *C `sm:"auto"`
}

func (d *D)Print()  {
	log.Println(d.Name, "[c]", d.C)

}
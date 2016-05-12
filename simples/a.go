package simples

import "log"

type A struct {
	Name     string
	D        *D `sm:"auto"`
	Password string `sm:"$.d.Password"`
}

func (a *A)Init()  {
	log.Println("hi init")
}
func (a *A)Print() {
	log.Println(a.Name,a.Password, "[d]", a.D)
}

type C struct {
	Name string
	A    *A `sm:"auto"`
}

func (c *C)Print() {
	log.Println(c.Name, "[a]", c.A)
}

type D struct {
	Name     string
	Password string `sm:"#.postgres.user"`
	C        *C `sm:"auto"`
}

func (d *D)Print() {
	log.Println(d.Name, "[c]", d.C)
}
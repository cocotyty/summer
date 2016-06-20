package main

import (
	"log"
	"os"
	"github.com/cocotyty/summer"
)

func main() {
	summer.Put(&P{})
	summer.Put(&PrinterProvider{})
	summer.Add("P",&PrinterProvider{prefix:"[p!]"})
	summer.Start()
}
type Printer interface {
	Println(args...interface{})
}

type P struct {
	Printer Printer `sm:"@.*"`
	Pr      Printer `sm:"@.P"`
	P       Printer `sm:"@.*"`
}

func (p *P)Ready() {
	p.Printer.Println("Printer is ready")
	p.Pr.Println("Printer is ready")
	p.P.Println("P is ready")
}

type PrinterProvider struct {
	prefix string
	logger *log.Logger
}

func (pp *PrinterProvider)Init() {
	pp.logger = log.New(os.Stderr, pp.prefix, log.LstdFlags)
}
func (pp *PrinterProvider)Provide() interface{} {
	return pp.logger
}
package summer

// a stone is a go stone as you know
type Stone interface{}

// a stone can init
type Init interface {
	Init()
}

// a stone can ready
type Ready interface {
	Ready()
}

// a stone can Destroy
type Destroy interface {
	Destroy()
}
type Provider interface {
	Provide() interface{}
}

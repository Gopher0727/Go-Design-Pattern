package main

type Target interface {
	Request() string
}

type Adaptee struct{}

func (a *Adaptee) SpecificRequest() string {
	return "被适配者的特定行为"
}

type Adapter struct {
	adaptee *Adaptee
}

func NewAdapter(adaptee *Adaptee) *Adapter {
	return &Adapter{adaptee: adaptee}
}

func (a *Adapter) Request() string {
	return "适配器: " + a.adaptee.SpecificRequest()
}

func main() {
	adaptee := &Adaptee{}
	adapter := NewAdapter(adaptee)

	callClient(adapter)
}

func callClient(t Target) {
	println(t.Request())
}

package main

import "fmt"

type ProtoType interface {
	Clone() ProtoType
}

func (p *Person) Clone() ProtoType {
	// 深拷贝
	var addrCopy *Address
	if p.Addr != nil {
		addrCopy = &Address{
			City:   p.Addr.City,
			Street: p.Addr.Street,
		}
	}
	return &Person{
		Name: p.Name,
		Age:  p.Age,
		Addr: addrCopy,
	}
}

type Address struct {
	City   string
	Street string
}

type Person struct {
	Name string
	Age  int
	Addr *Address
}

func (p *Person) String() string {
	addr := "<nil>"
	if p.Addr != nil {
		addr = fmt.Sprintf("%s, %s", p.Addr.City, p.Addr.Street)
	}
	return fmt.Sprintf("Person {Name: %s, Age: %d, Addr: %s}", p.Name, p.Age, addr)
}

// 注册表
type Registry map[string]ProtoType

func (r Registry) Register(name string, p ProtoType) {
	r[name] = p
}

func (r Registry) Get(name string) ProtoType {
	if p, ok := r[name]; ok {
		return p.Clone()
	}
	return nil
}

func main() {
	reg := Registry{}

	p := &Person{
		Name: "Alice",
		Age:  30,
		Addr: &Address{
			City:   "New York",
			Street: "5th Avenue",
		},
	}
	reg.Register("alice-prototype", p)

	np := reg.Get("alice-prototype").(*Person)
	np.Name = "Bob"
	np.Age = 25

	nnp := reg.Get("alice-prototype").(*Person)
	nnp.Name = "Charlie"
	nnp.Age = 28
	nnp.Addr.City = "Chicago"
	nnp.Addr.Street = "Michigan Avenue"

	fmt.Println(p)   // Original
	fmt.Println(np)  // Cloned and modified
	fmt.Println(nnp) // Cloned and modified
}

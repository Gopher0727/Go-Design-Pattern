package main

type House struct {
	Foundation string
	Structure  string
	Roof       string
}

func (h *House) String() string {
	return "[Foundation: " + h.Foundation + ", Structure: " + h.Structure + ", Roof: " + h.Roof + "]"
}

type HouseBuilder interface {
	BuildFoundation()
	BuildStructure()
	BuildRoof()

	GetHouse() *House

	Reset()
}

// 木结构房屋
type WoodenHouseBuilder struct {
	house *House
}

func NewWoodenHouseBuilder() *WoodenHouseBuilder {
	b := &WoodenHouseBuilder{}
	b.Reset()
	return b
}

func (b *WoodenHouseBuilder) Reset() {
	b.house = &House{}
}

func (b *WoodenHouseBuilder) BuildFoundation() {
	b.house.Foundation = "Wooden Foundation"
}

func (b *WoodenHouseBuilder) BuildStructure() {
	b.house.Structure = "Wooden Structure"
}

func (b *WoodenHouseBuilder) BuildRoof() {
	b.house.Roof = "Wooden Roof"
}

func (b *WoodenHouseBuilder) GetHouse() *House {
	house := b.house
	b.Reset()
	return house
}

// 砖结构房屋
type BrickHouseBuilder struct {
	house *House
}

func NewBrickHouseBuilder() *BrickHouseBuilder {
	b := &BrickHouseBuilder{}
	b.Reset()
	return b
}

func (b *BrickHouseBuilder) Reset() {
	b.house = &House{}
}

func (b *BrickHouseBuilder) BuildFoundation() {
	b.house.Foundation = "Brick Foundation"
}

func (b *BrickHouseBuilder) BuildStructure() {
	b.house.Structure = "Brick Structure"
}

func (b *BrickHouseBuilder) BuildRoof() {
	b.house.Roof = "Brick Roof"
}

func (b *BrickHouseBuilder) GetHouse() *House {
	house := b.house
	b.Reset()
	return house
}

// Director: 定义构建顺序
type Director struct {
	builder HouseBuilder
}

func NewDirector(b HouseBuilder) *Director {
	return &Director{builder: b}
}

func (d *Director) SetBuilder(b HouseBuilder) {
	d.builder = b
}

func (d *Director) Construct() {
	d.builder.BuildFoundation()
	d.builder.BuildStructure()
	d.builder.BuildRoof()
}

func main() {
	// 使用木结构房屋建造者
	woodenBuilder := NewWoodenHouseBuilder()
	director := NewDirector(woodenBuilder)
	director.Construct()
	woodenHouse := woodenBuilder.GetHouse()
	println("Wooden House:", woodenHouse.String())

	// 使用砖结构房屋建造者
	brickBuilder := NewBrickHouseBuilder()
	director.SetBuilder(brickBuilder)
	director.Construct()
	brickHouse := brickBuilder.GetHouse()
	println("Brick House:", brickHouse.String())
}
package pattern

import "fmt"

/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern
*/

/*
Строитель — это порождающий паттерн проектирования, который позволяет
создавать сложные объекты пошагово. Строитель даёт возможность
использовать один и тот же код строительства для получения разных
представлений объектов.

Строитель позволяет производить различные продукты, используя
один и тот же процесс строительства.

Применение:
- Когда вы хотите избавиться от "телескопического конструктора".
- Когда ваш код должен создавать разные представления какого-то объекта.
- Когда вам нужно собирать сложные составные объекты.

Плюсы:
- Позволяет создавать продукты пошагово.
- Позволяет использовать один и тот же код для создания различных продуктов.
- Изолирует сложный код сборки продукта от его основной бизнес-логики.

Минусы:
- Усложняет код программы из-за введения дополнительных классов.
- Клиент будет привязан к конкретным классам строителей, так как в
интерфейсе директора может не быть метода получения результата.
*/

type builderType string

const (
	MSIBuilderType  builderType = "MSI"
	AsusBuilderType builderType = "Asus"
)

type Builder interface {
	SetBrand()
	SetMotherboard()
	SetCPU()
	SetGPU()
	SetMemory()
	SetDisk()
	SetMonitor()
	GetComputer() Computer
}

func GetBuilder(t builderType) Builder {
	switch t {
	case MSIBuilderType:
		return &MSIBuilder{}
	case AsusBuilderType:
		return &AsusBuilder{}
	default:
		return nil
	}
}

type Computer struct {
	Brand       string
	Motherboard string
	CPU         string
	GPU         string
	Memory      string
	Disk        string
	Monitor     string
}

func (c *Computer) Show() {
	fmt.Printf("Brand: %s\n", c.Brand)
	fmt.Printf("Motherboard: %s\n", c.Motherboard)
	fmt.Printf("CPU: %s\n", c.CPU)
	fmt.Printf("GPU: %s\n", c.GPU)
	fmt.Printf("Memory: %s\n", c.Memory)
	fmt.Printf("Disk: %s\n", c.Disk)
	fmt.Printf("Monitor: %s\n", c.Monitor)
}

type MSIBuilder struct {
	Brand       string
	Motherboard string
	CPU         string
	GPU         string
	Memory      string
	Disk        string
	Monitor     string
}

func (b *MSIBuilder) SetBrand() {
	b.Brand = "MSI"
}
func (b *MSIBuilder) SetMotherboard() {
	b.Motherboard = "MSI MPG Z690"
}

func (b *MSIBuilder) SetCPU() {
	b.CPU = "i7-10750H"
}

func (b *MSIBuilder) SetGPU() {
	b.GPU = "GTX 1650"
}

func (b *MSIBuilder) SetMemory() {
	b.Memory = "16GB"
}

func (b *MSIBuilder) SetDisk() {
	b.Disk = "512GB"
}

func (b *MSIBuilder) SetMonitor() {
	b.Monitor = "MSI G7-H144v3"
}

func (b *MSIBuilder) GetComputer() Computer {
	return Computer{
		Brand:       b.Brand,
		Motherboard: b.Motherboard,
		CPU:         b.CPU,
		GPU:         b.GPU,
		Memory:      b.Memory,
		Disk:        b.Disk,
		Monitor:     b.Monitor,
	}
}

type AsusBuilder struct {
	Brand       string
	Motherboard string
	CPU         string
	GPU         string
	Memory      string
	Disk        string
	Monitor     string
}

func (b *AsusBuilder) SetBrand() {
	b.Brand = "Asus"
}
func (b *AsusBuilder) SetMotherboard() {
	b.Motherboard = "Asus ROG Strix B550-F"
}

func (b *AsusBuilder) SetCPU() {
	b.CPU = "Ryzen 5 3600"
}

func (b *AsusBuilder) SetGPU() {
	b.GPU = "RTX 2060"
}

func (b *AsusBuilder) SetMemory() {
	b.Memory = "8GB"
}

func (b *AsusBuilder) SetDisk() {
	b.Disk = "1024GB"
}

func (b *AsusBuilder) SetMonitor() {
	b.Monitor = "Asus ROG G-H120m1"
}

func (b *AsusBuilder) GetComputer() Computer {
	return Computer{
		Brand:       b.Brand,
		Motherboard: b.Motherboard,
		CPU:         b.CPU,
		GPU:         b.GPU,
		Memory:      b.Memory,
		Disk:        b.Disk,
		Monitor:     b.Monitor,
	}
}

type Director struct {
	Builder Builder
}

func NewDirector(b Builder) *Director {
	return &Director{
		Builder: b,
	}
}

func (d *Director) SetBuilder(b Builder) {
	d.Builder = b
}

// Основная функция по строительству объекта,
// в ней можно изменять порядок строительства
func (d *Director) CreateComputer() Computer {
	d.Builder.SetBrand()
	d.Builder.SetMotherboard()
	d.Builder.SetCPU()
	d.Builder.SetGPU()
	d.Builder.SetMemory()
	d.Builder.SetDisk()
	d.Builder.SetMonitor()
	return d.Builder.GetComputer()
}

func main() {
	msiBuilder := GetBuilder(MSIBuilderType)
	asusBuilder := GetBuilder(AsusBuilderType)

	director := NewDirector(msiBuilder)
	msiComputer := director.CreateComputer()
	msiComputer.Show()

	fmt.Println()

	director.SetBuilder(asusBuilder)
	asusComputer := director.CreateComputer()
	asusComputer.Show()

}

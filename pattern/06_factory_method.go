package pattern

import "fmt"

/*
	Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern
*/

/*
Фабричный метод — это порождающий паттерн проектирования, который определяет
общий интерфейс для создания объектов в суперклассе,
позволяя подклассам изменять тип создаваемых объектов.

Фабричный метод задаёт метод, который следует использовать вместо
вызова оператора `new` для создания объектов-продуктов. Подклассы
могут переопределить этот метод, чтобы изменять тип создаваемых продуктов.

Применение:
- Когда заранее неизвестны типы и зависимости объектов, с которыми
должен работать ваш код.
- Когда вы хотите дать возможность пользователям расширять части
вашего фреймворка или библиотеки.

Плюсы:
- Избавляет класс от привязки к конкретным классам продуктов.
- Выделяет код производства продуктов в одно место, упрощая поддержку кода.
- Упрощает добавление новых продуктов в программу.
- Реализует принцип открытости/закрытости.

Минусы:
- Может привести к созданию больших параллельных иерархий классов,
так как для каждого продукта надо создавать свой подкласс создателя.
*/

type NonPlayerCharacter interface {
	SetName(name string)
	SetLevel(level int)
	Name() string
	Level() int
}

type NPC struct {
	name  string
	level int
}

func (npc *NPC) SetName(name string) {
	npc.name = name
}

func (npc *NPC) SetLevel(level int) {
	npc.level = level
}

func (npc *NPC) Name() string {
	return npc.name
}

func (npc *NPC) Level() int {
	return npc.level
}

type Human struct {
	NPC
}

func NewHuman(name string, level int) *Human {
	return &Human{
		NPC: NPC{
			name:  name,
			level: level,
		},
	}
}

type Dwarf struct {
	NPC
}

func NewDwarf(name string, level int) *Dwarf {
	return &Dwarf{
		NPC: NPC{
			name:  name,
			level: level,
		},
	}
}

type NPCType string

const (
	HumanNPC NPCType = "Human"
	DwarfNPC NPCType = "Dwarf"
)

func NewNPC(typeNPC NPCType, name string, level int) (NonPlayerCharacter, error) {
	switch typeNPC {
	case HumanNPC:
		return NewHuman(name, level), nil
	case DwarfNPC:
		return NewDwarf(name, level), nil
	}
	return nil, fmt.Errorf("Неверный тип NPC")
}

func main() {
	human, err := NewNPC(HumanNPC, "Boromir", 23)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("npc: %s\tlevel: %d\n", human.Name(), human.Level())

	dwarf, err := NewNPC(DwarfNPC, "Gimli", 47)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("npc: %s\tlevel: %d\n", dwarf.Name(), dwarf.Level())

}

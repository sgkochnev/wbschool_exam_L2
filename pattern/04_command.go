package pattern

import "fmt"

/*
	Реализовать паттерн «команда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern
*/

/*
Команда — это поведенческий паттерн проектирования, который превращает
запросы в объекты, позволяя передавать их как аргументы при вызове методов,
ставить запросы в очередь, логировать их, а также поддерживать отмену операций.

Команда позволяет передавать запросы в очередь, логировать
их и поддерживать отмену операций.
Паттерн Команда устраняет прямую зависимость между объектами,
вызывающими операции, и объектами, выполняющими их.
Команда позволяет реализовать простую отмену и повтор операций,
а также отложенный запуск операций.

Применение:
- Когда вы хотите параметризировать объекты выполняемым действием.
- Когда вы хотите ставить операции в очередь, выполнять их по расписанию
или передавать по сети.
- Когда вам нужна операция отмены.

Плюсы:
- Убирает прямую зависимость между объектами, вызывающими операции,
и объектами, которые их непосредственно выполняют.
- Позволяет реализовать простую отмену операций.
- Позволяет реализовать отложенный запуск операций.
- Позволяет собирать сложные команды из простых.
- Реализует принцип открытости/закрытости.

Минусы:
- Усложняет код программы из-за введения множества дополнительных классов.
*/

type Command interface {
	Execute()
}

type Button struct {
	command Command
}

func (b *Button) Press() {
	b.command.Execute()
}

type Device interface {
	On()
	Off()
}

type OnCommand struct {
	device Device
}

func (c *OnCommand) Execute() {
	c.device.On()
}

type OffCommand struct {
	device Device
}

func (c *OffCommand) Execute() {
	c.device.Off()
}

type TV struct {
	isRunning bool
}

func (t *TV) On() {
	t.isRunning = true
	fmt.Println("Turning TV on")
}

func (t *TV) Off() {
	t.isRunning = false
	fmt.Println("Turning TV off")
}

func main() {
	tv := &TV{}
	onCommand := &OnCommand{device: tv}
	offCommand := &OffCommand{device: tv}

	onButton := &Button{command: onCommand}
	onButton.Press()

	offButton := &Button{command: offCommand}
	offButton.Press()
}

package pattern

import "fmt"

/*
	Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern
*/

/*
Состояние — это поведенческий паттерн проектирования, который позволяет
объектам менять поведение в зависимости от своего состояния. Извне
создаётся впечатление, что изменился класс объекта.

Применимость:
- Когда у вас есть объект, поведение которого кардинально меняется в зависимости от
внутреннего состояния, причём типов состояний много, и их код часто меняется.
- Когда код класса содержит множество больших, похожих друг на друга,
условных операторов,
которые выбирают поведения в зависимости от текущих значений полей класса.
- Когда вы сознательно используете табличную машину состояний,
построенную на условных операторах, но вынуждены мириться
с дублированием кода для похожих состояний и переходов.

Плюсы:
- Избавляет от множества больших условных операторов машины состояний.
- Концентрирует код в одном месте, связанный с определенным состоянием.
- Упрощает код контекста.

Минусы:
- Может непосредственно усложнить код, если состояний мало и они редко меняются.
*/
type State interface {
	Action()
	String() string
}

type AcceptedState struct {
	parcel *Parcel
}

func (a *AcceptedState) Action() {
	a.parcel.SetState(a)
}

func (a *AcceptedState) String() string {
	return "Посылка принята в обработку."
}

type InTransitState struct {
	parcel *Parcel
}

func (it *InTransitState) Action() {
	it.parcel.SetState(it)
}

func (it *InTransitState) String() string {
	return "Посылка в пути."
}

type DeliveredState struct {
	parcel *Parcel
}

func (d *DeliveredState) Action() {
	d.parcel.SetState(d)
}

func (d *DeliveredState) String() string {
	return "Посылка доставлена."
}

type Parcel struct {
	state State
}

func (p *Parcel) SetState(state State) {
	p.state = state
}

func (p *Parcel) State() string {
	return p.state.String()
}

func main() {
	parcel := &Parcel{}

	accepted := &AcceptedState{parcel}
	accepted.Action()
	fmt.Println(parcel.State())

	inTransit := &InTransitState{parcel}
	inTransit.Action()
	fmt.Println(parcel.State())

	delivered := &DeliveredState{parcel}
	delivered.Action()
	fmt.Println(parcel.State())
}

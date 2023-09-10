package pattern

import "fmt"

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern
*/

/*
Стратегия — это поведенческий паттерн проектирования, который определяет
семейство схожих алгоритмов и помещает каждый из них в собственный класс,
после чего алгоритмы можно взаимозаменять прямо во время исполнения программы.

Применение:
- Когда вам нужно использовать разные вариации какого-то алгоритма внутри
одного объекта.
- Когда у вас есть множество похожих классов, отличающихся только
некоторым поведением.
- Когда вы не хотите раскрывать детали реализации алгоритмов для
других классов.
- Когда различные вариации алгоритмов реализованы в виде развесистого
условного оператора. Каждая ветка такого оператора представляет собой
вариацию алгоритма.

Плюсы:
- Горячая замена алгоритмов на лету.
- Изолирует  код и данные алгоритма от остальных классоав.
- Уход от наследования к делигированию.
- Реализует принцип открытости/закрытости.

Минусы:
- Усложняет программу за счет дополнительных классов.
- Клиент должен знать, в чем состоит разница между стратегиями,
чтобы выбрать подходящую.
*/

type Strategy interface {
	Execute(float64, float64) float64
}

type StrategyAdd struct {
	a, b float64
}

func (s *StrategyAdd) Execute(a, b float64) float64 {
	return a + b
}

type StrategySub struct {
	a, b float64
}

func (s *StrategySub) Execute(a, b float64) float64 {
	return a - b
}

type StrategyMul struct {
	a, b float64
}

func (s *StrategyMul) Execute(a, b float64) float64 {
	return a * b
}

type StrategyDiv struct {
	a, b float64
}

func (s *StrategyDiv) Execute(a, b float64) float64 {
	return a / b
}

type Context struct {
	strategy Strategy
}

func (c *Context) SetStrategy(strategy Strategy) {
	c.strategy = strategy
}

func (c *Context) Execute(a, b float64) float64 {
	return c.strategy.Execute(a, b)
}

func main() {
	c := &Context{}
	c.SetStrategy(&StrategyAdd{})
	fmt.Println(c.Execute(5, 3))

	c.SetStrategy(&StrategySub{})
	fmt.Println(c.Execute(5, 3))

	c.SetStrategy(&StrategyMul{})
	fmt.Println(c.Execute(5, 3))

	c.SetStrategy(&StrategyDiv{})
	fmt.Println(c.Execute(5, 3))
}

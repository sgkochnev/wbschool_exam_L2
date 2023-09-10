package pattern

import (
	"fmt"
	"math"
)

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern
*/

/*
Посетитель — это поведенческий паттерн, который позволяет добавить
новую операцию для целой иерархии классов, не изменяя код этих классов.

Паттерн Посетитель позволяет вам добавлять поведение в структуру без ее изменения.
Посетитель позволяет применять одну и ту же операцию к объектам различных классов.

Примнение:
- Когда вам нужно выполнить какую-то операцию над всеми элементами сложной
структуры объектов.
- Когда над объектами сложной структуры надо выполнить некоторые не связанные
между собой операции, но вы не хотите "засорять" классы такими операциями.
- Когда новое поведение имеет смысл только для некоторых классов из
существующей иерархии.

Плюсы:
- Упрощает добавление операций, работающих со сложными структурами объектов.
- Объединяет родственные операции в одном классе.const
- Посетитель может накапливать состояние при обходе структуры элементов.const

Минусы:
- Паттерн не оправдан, если иерархия элементов часто меняется.const
- Может привести к нарушению инкапсуляции элементов.
*/

type Visitor interface {
	VisitForRectangle(*Rectangle)
	VisitForSquare(*Square)
	VisitForCircle(*Circle)
}

type Shape interface {
	GetType() string
	Accept(Visitor)
}

type Square struct {
	side float64
}

func (s *Square) GetType() string {
	return "Square"
}

func (s *Square) Accept(v Visitor) {
	v.VisitForSquare(s)
}

type Rectangle struct {
	width  float64
	height float64
}

func (r *Rectangle) GetType() string {
	return "Rectangle"
}

func (r *Rectangle) Accept(v Visitor) {
	v.VisitForRectangle(r)
}

type Circle struct {
	radius float64
}

func (c *Circle) GetType() string {
	return "Circle"
}

func (c *Circle) Accept(v Visitor) {
	v.VisitForCircle(c)
}

type AreaCalculator struct {
	area float64
}

func NewAreaCalculator() *AreaCalculator {
	return &AreaCalculator{}
}

func (ac *AreaCalculator) VisitForRectangle(r *Rectangle) {
	ac.area = r.width * r.height
	fmt.Printf("Area fot rectangle: %0.3f\n", ac.area)
}

func (ac *AreaCalculator) VisitForSquare(s *Square) {
	ac.area = s.side * s.side
	fmt.Printf("Area for square: %0.3f\n", ac.area)
}

func (ac *AreaCalculator) VisitForCircle(c *Circle) {
	ac.area = c.radius * c.radius * math.Pi
	fmt.Printf("Area for circle: %0.3f\n", ac.area)
}

func (ac *AreaCalculator) GetArea() float64 {
	return ac.area
}

func main() {
	square := &Square{side: 5}
	rectangle := &Rectangle{width: 10, height: 5}
	circle := &Circle{radius: 5}

	areaCalculator := NewAreaCalculator()

	rectangle.Accept(areaCalculator)
	fmt.Println(areaCalculator.GetArea())

	square.Accept(areaCalculator)
	fmt.Println(areaCalculator.GetArea())

	circle.Accept(areaCalculator)
	fmt.Println(areaCalculator.GetArea())

}

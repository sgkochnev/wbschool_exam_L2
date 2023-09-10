package pattern

import "fmt"

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/

/*
Цепочка обязанностей — это поведенческий паттерн проектирования, который
позволяет передавать запросы последовательно по цепочке обработчиков.
Каждый последующий обработчик решает, может ли он обработать запрос сам
и стоит ли передавать запрос дальше по цепи.

Избавляет от жёсткой привязки отправителя запроса к его получателю,
позволяя выстраивать цепь из различных обработчиков динамически.

Применение:
- Когда программа должна обрабатывать разнообразные запросы несколькими
способами, но заранее неизвестно, какие конкретно запросы будут происходить
и какие обработчики для них понадабятся.
- Когда важно, чтобы обработчик выполнялись один за другим в строгом порядке.
- Когда набор объектов, способных обработать запрос, должен задаваться
динамически.

Плюсы:
- Уменьшает зависимость между клиентом и обработчиками.
- Реализует принцип единственной оюязанности.
- Реализует принцип открытости/закрытости.

Минусы:
- Запрос может остаться никем не обработанным.
*/

type Service interface {
	Execute(*Data)
	SetNext(Service)
}

type Data struct {
	GetSource    bool
	UpdateSource bool
}

type Device struct {
	Name string
	Next Service
}

func (d *Device) Execute(data *Data) {
	if !data.GetSource {
		fmt.Printf("Get data from device %s.\n", d.Name)
		data.GetSource = true
	} else {
		fmt.Printf("Data from device %s alrady get.\n", d.Name)
	}

	if d.Next != nil {
		d.Next.Execute(data)
	}
}

func (d *Device) SetNext(service Service) {
	d.Next = service
}

type UpdateDataService struct {
	Name string
	Next Service
}

func (uds *UpdateDataService) Execute(data *Data) {
	if !data.UpdateSource {
		fmt.Printf("Update data from service %s.\n", uds.Name)
		data.UpdateSource = true
	} else {
		fmt.Printf("Data in service %s alrady update.\n", uds.Name)
	}

	if uds.Next != nil {
		uds.Next.Execute(data)
	}
}

func (uds *UpdateDataService) SetNext(service Service) {
	uds.Next = service
}

type SaveDataService struct {
	Next Service
}

func (sds *SaveDataService) Execute(data *Data) {
	if !data.UpdateSource {
		fmt.Println("Data not update.")
	} else {
		fmt.Println("Data save.")
	}

	if sds.Next != nil {
		sds.Next.Execute(data)
	}
}

func (sds *SaveDataService) SetNext(service Service) {
	sds.Next = service
}

func main() {
	device := &Device{Name: "Device1"}
	updateService := &UpdateDataService{Name: "DataService1"}
	saveService := &SaveDataService{}

	device.SetNext(updateService)
	updateService.SetNext(saveService)

	data := &Data{GetSource: true}
	device.Execute(data)
}

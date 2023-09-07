Что выведет программа? Объяснить вывод программы.

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)

	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}

		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v := <-a:
				c <- v
			case v := <-b:
				c <- v
			}
		}
	}()
	return c
}

func main() {

	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4 ,6, 8)
	c := merge(a, b )
	for v := range c {
		fmt.Println(v)
	}
}
```

Ответ:
```
В случайном порядке выведутся цифры от 1 до 8,
потом будет бесконечный поток значений по умолчанию (для `int` это 0). 

Так произойдет потому, что каналы `a` и `b` будут закрыты после того как 
обработают все переденные в функцию значения. Канал `c` будет читать 
из закрытых каналов `a` и `b`.
Чтение из закрытого канала дает значение по умолчанию для типа канала. 

```

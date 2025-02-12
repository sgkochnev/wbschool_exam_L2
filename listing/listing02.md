Что выведет программа? Объяснить вывод программы. Объяснить как работают defer’ы и их порядок вызовов.

```go
package main

import (
	"fmt"
)


func test() (x int) {
	defer func() {
		x++
	}()
	x = 1
	return
}


func anotherTest() int {
	var x int
	defer func() {
		x++
	}()
	x = 1
	return x
}


func main() {
	fmt.Println(test()) 		// 2
	fmt.Println(anotherTest()) 	// 1
}
```

Ответ:
```
2
1

Отложенные функции (defer) выполняются после завершения функции, в которой 
они были объявлены, но перед тем, как функция вернет значение.
Если вызвать несколько defer'ов, они выполняются в обратном порядке.

В случае наличия в функции именованных возвращаемых значений 
функция в defer может как читать,
так и модифицировать эти именованные возвращаемые значения.
Если функция в defer модифицирует именованное возвращаемое значение,
то будет возвращено именно это модифицированное значение.

```

Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
error

Из функции `test` возвращается типизированная ощибка.
Результатом сравнения типизированной ошибки с `nil` всегда будет `false`.

```

Более подробно расписано в [listing03.md](./listing03.md).

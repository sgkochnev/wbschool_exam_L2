package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

/*
=== Базовая задача ===

Создать программу печатающую точное время с использованием NTP библиотеки.Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/

const address = "0.beevik-ntp.pool.ntp.org"

// Time returns the current, corrected local time using information returned from the remote NTP server. On error, Time returns the uncorrected local system time.
func Time() (time.Time, error) {
	return ntp.Time(address)
}

func main() {
	t, err := Time()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(t.Format("2006-01-02 15:04:05.000"))
}

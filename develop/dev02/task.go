package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// ErrInvalidtString is returned when the input string is invalid.
var ErrInvalidtString = errors.New("invalid string")

const backslash = '\\'

// Unpack takes a string and unpacks it according to the specified rules.
//
// For example:
// "a4bc2d5e" => "aaaabccddddde"
// "qwe\45" => qwe44444
// "qwe45" - invalid string
//
// The function accepts a single parameter:
// - str: a string that needs to be unpacked.
//
// It returns two values:
// - string: the unpacked version of the input string.
// - error: an error if the input string is invalid.
func Unpack(str string) (string, error) {

	var (
		prevCahr        *rune
		prevCharIsDigit bool
		escaped         bool
	)

	builder := &strings.Builder{}

	for i, v := range str {
		v := v
		switch {
		case escaped: // current character is escaped
			if !unicode.IsDigit(v) && v != backslash {
				return "", ErrInvalidtString
			}
			prevCahr = &v
			escaped = false

		case v == backslash: // escape character
			if prevCahr != nil {
				builder.WriteRune(*prevCahr)
				prevCahr = nil
			}
			escaped = true

		case unicode.IsDigit(v): // current character is digit
			if prevCharIsDigit || i == 0 {
				return "", ErrInvalidtString
			}

			prevCharIsDigit = i != len(str)-1 // true if not last

			count := int(v - '0') // convert digit form `rune` to `int`
			if prevCahr != nil {
				builder.WriteString(strings.Repeat(string(*prevCahr), count))
				prevCahr = nil
			}

		default: // current character is not digit and not escaped
			if prevCahr != nil {
				builder.WriteRune(*prevCahr)
			}
			prevCharIsDigit = false
			prevCahr = &v
		}
	}

	// last character processing
	if prevCharIsDigit || escaped {
		return "", ErrInvalidtString
	}

	if prevCahr != nil {
		builder.WriteRune(*prevCahr)
	}

	return builder.String(), nil
}

func main() {
	fmt.Println(Unpack(`a4bc2d5e`))
	fmt.Println(Unpack(`abcd`))
	fmt.Println(Unpack(`45`))
	fmt.Println(Unpack(``))
	fmt.Println(Unpack(`qwe\4\5`))
	fmt.Println(Unpack(`qwe\45`))
	fmt.Println(Unpack(`qwe\\5`))
	fmt.Println(Unpack(`asd\a`))
}

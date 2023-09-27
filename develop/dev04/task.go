package main

import (
	"fmt"
	"slices"
	"strings"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func findAnagrams(words []string) map[string][]string {
	keys := make(map[string]string)
	vals := make(map[string]map[string]struct{})

	for _, word := range words {
		word = strings.ToLower(word)
		k := []rune(word)
		slices.Sort(k)

		key, ok := keys[string(k)]
		if !ok {
			keys[string(k)] = word
			key = word
		}

		if _, ok := vals[key]; !ok {
			vals[key] = make(map[string]struct{})
		}
		vals[key][word] = struct{}{}
	}

	res := make(map[string][]string)

	for key, val := range vals {
		if len(val) > 1 {

			sl := make([]string, 0, len(val))
			for k := range val {
				sl = append(sl, k)
			}
			slices.Sort(sl)

			res[key] = sl
		}
	}

	return res
}

func main() {
	words := []string{"тяпка", "пятак", "Пятка", "Пятка", "слиток", "листок", "Столик", "свисток", "рулетка"}

	m := findAnagrams(words)
	for k, v := range m {
		fmt.Println(k, v)
	}
}

package main

import (
	"bufio"
	"cmp"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {

	k := flag.Int("k", 0, "колонка для сортировки") // +
	sep := flag.String("sep", " ", "символ разделителя на колони (по умолчанию пробел)")
	n := flag.Bool("n", false, "сортировать по числовому значению") // +
	r := flag.Bool("r", false, "сортировать в обратном порядке")    // +
	u := flag.Bool("u", false, "не выводить повторяющиеся строки")  // +
	M := flag.Bool("M", false, "сортировать по названию месяца")    // +
	b := flag.Bool("b", false, "игнорировать хвостовые пробелы")    // +
	c := flag.Bool("c", false, "проверять отсортированы ли данные") // +

	flag.Parse()

	log.Println(*sep)

	flags := []bool{*n, *M, *c}

	count := 0

	for _, v := range flags {
		if v {
			count++
		}
	}

	if count > 1 {
		log.Fatalln("Можно использовать только один из ключей -n, -M, -c")
	}

	if *k == 0 {
		k = nil
	}

	args := flag.Args()
	if len(args) == 0 {
		log.Fatalln("Не указано имя файла")
	}

	data, err := NewDataFromFile(args[0], *b)
	if err != nil {
		fmt.Println(err)
	}

	out := os.Stdout
	if len(args) >= 2 {
		out, err = os.OpenFile(args[1], os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
	}

	switch {
	case *c:
		var isSorted bool
		if *n || *M {
			isSorted = data.IsSorted(cmpNumbersAndMonths(k, *M, *sep, *r))
		} else {
			isSorted = data.IsSorted(cmpDefault(k, *sep, *r))
		}
		fmt.Fprintln(out, isSorted)
		return

	case *n, *M:
		data.Sort(cmpNumbersAndMonths(k, *M, *sep, *r))
	default:
		data.Sort(cmpDefault(k, *sep, *r))
	}

	if *u {
		data.DeleteRepeatedLines(k, *sep)
	}

	writeData(out, data)
}

func writeData(out *os.File, data *Data) {
	writer := bufio.NewWriter(out)
	for _, s := range data.text {
		_, err := writer.WriteString(s)
		if err != nil {
			log.Println(err)
		}
		_, _ = writer.WriteRune('\n')
	}
	writer.Flush()
}

type Data struct {
	text []string
}

func (d *Data) Sort(cmp func(a, b string) int) {
	slices.SortFunc(d.text, cmp)
}

func (d *Data) IsSorted(cmp func(a, b string) int) bool {
	return slices.IsSortedFunc(d.text, cmp)
}

func (d *Data) DeleteRepeatedLines(col *int, sep string) {
	var s1, s2 string

	if col != nil {
		s2s := strings.Split(d.text[0], sep)
		if len(s2s) <= *col-1 {
			log.Fatalf("Не могу разбить строку на %d или более колонок", *col)
		}
		s2 = s2s[*col-1]
	} else {
		s2 = d.text[0]
	}

	j := 1
	for i := 1; i < len(d.text); i++ {
		if col != nil {
			s1s := strings.Split(d.text[i], sep)
			if len(s1s) <= *col-1 {
				log.Fatalf("Не могу разбить строку на %d или более колонок", *col)
			}
			s1 = s1s[*col-1]
		} else {
			s1 = d.text[i]
		}

		if s1 != s2 {
			d.text[j] = d.text[i]
			j++
		}
		s2 = s1
	}
	d.text = d.text[:j]
}

// NewDataFromFile reads a file and returns a new Data object.
func NewDataFromFile(filepath string, trimSpace bool) (*Data, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	text := make([]string, 0)

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if trimSpace {
			text = append(text, strings.TrimRightFunc(string(line), unicode.IsSpace))
		} else {
			text = append(text, string(line))
		}
	}

	return &Data{text: text}, nil
}

func mustConvertMonth(s string) string {
	d, err := time.Parse("January", s)
	if err != nil {
		log.Fatalln(err)
	}
	return d.Format("01")
}

func cmpNumbersAndMonths(col *int, month bool, sep string, reverse bool) func(a, b string) int {
	return func(a, b string) int {

		if col != nil {
			as := strings.Split(a, sep)
			if len(as) <= *col-1 {
				log.Fatalf("Не могу разбить на %d или более колонок", *col)
			}
			bs := strings.Split(b, sep)
			if len(bs) <= *col-1 {
				log.Fatalf("Не могу разбить на %d или более колонок", *col)
			}
			a = as[*col-1]
			b = bs[*col-1]
		}

		if month {
			a = mustConvertMonth(a)
			b = mustConvertMonth(b)
		}

		num1, err := strconv.ParseFloat(a, 64)
		if err != nil {
			log.Println(err)
		}
		num2, err := strconv.ParseFloat(b, 64)
		if err != nil {
			log.Println(err)
		}

		if reverse {
			return cmp.Compare(num2, num1)
		}
		return cmp.Compare(num1, num2)
	}
}

func cmpDefault(col *int, sep string, reverse bool) func(a, b string) int {
	return func(a, b string) int {
		if col != nil {
			as := strings.Split(a, sep)
			if len(as) <= *col-1 {
				log.Fatalf("Не могу разбить на %d или более колонок", *col)
			}
			bs := strings.Split(b, sep)
			if len(bs) <= *col-1 {
				log.Fatalf("Не могу разбить на %d или более колонок", *col)
			}
			a = as[*col-1]
			b = bs[*col-1]
		}

		if reverse {
			return cmp.Compare(b, a)
		}
		return cmp.Compare(a, b)
	}
}

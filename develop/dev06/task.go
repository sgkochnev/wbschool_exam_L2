package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type flags struct {
	fields    [][2]int
	delimiter string
	separated bool
}

var ErrNotNumber = errors.New("колонки должны быть в числовом формате (например: '1,2,4' или '1-4')")
var ErrInvalidInterval = errors.New("не корректный интервал (например: '1-3,6' или '-3,6')")

func main() {

	f := &flags{}

	flag.Func("f", "выбрать поля (колонки)", parseFields(f))
	flag.StringVar(&f.delimiter, "d", "\t", "использовать другой разделитель")
	flag.BoolVar(&f.separated, "s", false, "только строки с разделителем")

	flag.Parse()

	in := os.Stdin
	if flag.NArg() > 0 {
		var err error
		in, err = os.Open(flag.Arg(0))
		if err != nil {
			panic(err)
		}
		defer in.Close()
	}

	out := os.Stdout
	if flag.NArg() > 1 {
		var err error
		out, err = os.OpenFile(flag.Arg(1), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			panic(err)
		}
		defer out.Close()
	}

	cut(in, out, f)
}

func cut(r io.Reader, w io.Writer, f *flags) {
	reader := bufio.NewReader(r)
	writer := bufio.NewWriter(w)
	defer writer.Flush()

	d := []byte(f.delimiter)

	for {

		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		if f.separated && !bytes.Contains(line, d) {
			continue
		}

		cols := bytes.Split(line, d)

		l := len(cols)
		n := 0

		for _, v := range f.fields {
			if v[0] > l {
				break
			}

			for i := v[0]; i <= v[1] && i <= l; i++ {
				if n > 0 {
					_, _ = writer.Write(d)
				}
				k, _ := writer.Write(cols[i-1])
				n += k
			}

			if v[1] > l {
				break
			}
		}

		if n > 0 {
			_, _ = writer.WriteRune('\n')
		}
	}
}

func parseFields(f *flags) func(string) error {
	return func(s string) error {
		var err error

		fields := strings.Split(s, ",")
		for _, v := range fields {

			cols := strings.Split(v, "-")

			interval := [2]int{1, math.MaxInt} // начальное значение интервала

			if cols[0] != "" {
				interval[0], err = strconv.Atoi(cols[0])
				if err != nil {
					return ErrNotNumber
				}
				interval[0] = max(interval[0], 1)
			}

			if len(cols) > 2 {
				return fmt.Errorf("неожидаемый интервал %s : %w", v, ErrInvalidInterval)
			}

			if len(cols) == 2 {
				if cols[1] != "" {
					interval[1], err = strconv.Atoi(cols[1])
					if err != nil {
						return ErrNotNumber
					}
					interval[1] = max(interval[1], interval[0])
				}
			} else {
				interval[1] = interval[0]
			}

			f.fields = append(f.fields, interval)
		}

		f.fields = merge(f.fields)

		return nil
	}
}

func merge(intervals [][2]int) [][2]int {
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	res := make([][2]int, 0)

	for _, v := range intervals {
		if len(res) == 0 || res[len(res)-1][1] < v[0] {
			res = append(res, v)
		} else {
			res[len(res)-1][1] = max(res[len(res)-1][1], v[1])
		}
	}
	return res
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения +
-B - "before" печатать +N строк до совпадения +
-C - "context" (A+B) печатать ±N строк вокруг совпадения +
-c - "count" (количество строк) +
-i - "ignore-case" (игнорировать регистр) +
-v - "invert" (вместо совпадения, исключать) +
-F - "fixed", точное совпадение со строкой, не паттерн -
-n - "line num", печатать номер строки +

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type ringQueue[T any] struct {
	s, e int
	buf  []T
	l    int
}

func newRingQueue[T any](size int) *ringQueue[T] {
	return &ringQueue[T]{
		buf: make([]T, size),
	}
}

func (r *ringQueue[T]) push(v T) {
	if r.l == len(r.buf) {
		r.s = (r.s + 1) % len(r.buf)
		r.l--
	}
	r.buf[r.e] = v
	r.e = (r.e + 1) % len(r.buf)
	r.l++
}

func (r *ringQueue[T]) pop() (T, bool) {
	if r.l > 0 {
		v := r.buf[r.s]
		r.s = (r.s + 1) % len(r.buf)
		r.l--
		return v, true
	}
	var v T
	return v, false
}

func (r *ringQueue[T]) len() int {
	return r.l
}

func (r *ringQueue[T]) isEmpty() bool {
	return r.len() == 0
}

type flags struct {
	after, before, context                    int
	count, ignoreCase, invert, fixed, lineNum bool
	pattern                                   string
}

func pattern(str string, f *flags) string {
	prefix := ""
	if f.ignoreCase {
		prefix += "i" // (?i) - ignore case
	}
	if f.fixed {
		prefix += "sm"
		str = fmt.Sprintf("^\\Q%s\\E$", str) // ^\Q...\E$ - fixed
	}

	if prefix != "" {
		prefix = fmt.Sprintf("(?%s)", prefix)
		str = fmt.Sprintf("%s%s", prefix, str)
	} else {
		str = fmt.Sprint(str)
	}

	return str
}

type buf struct {
	line []byte
	num  int
}

func grep(r io.Reader, w io.Writer, f *flags) {
	writer := bufio.NewWriter(w)
	defer writer.Flush()

	reader := bufio.NewReader(r)

	re, err := regexp.Compile(pattern(f.pattern, f))
	if err != nil {
		panic(err)
	}

	f.after = max(f.after, f.context)
	f.before = max(f.before, f.context)

	bufBefore := newRingQueue[buf](f.before)

	countAfter := 0

	count := 0
	lineNum := 0

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		lineNum++

		match := re.Match(line)
		if f.invert {
			match = !match
		}

		if match {

			if f.count {
				count++
			} else {
				for !bufBefore.isEmpty() {
					b, _ := bufBefore.pop()
					if f.lineNum {
						_, _ = fmt.Fprintf(writer, "%d:", b.num)
					}
					_, _ = writer.Write(b.line)
					_, _ = writer.WriteRune('\n')
				}
				if f.lineNum {
					_, _ = fmt.Fprintf(writer, "%d:", lineNum)
				}
				_, _ = writer.Write(line)
				_, _ = writer.WriteRune('\n')
				countAfter = f.after
			}
		} else if !f.count {
			if countAfter > 0 {
				if f.lineNum {
					_, _ = fmt.Fprintf(writer, "%d:", lineNum)
				}
				_, _ = writer.Write(line)
				_, _ = writer.WriteRune('\n')
				countAfter--
				continue
			}

			if f.before > 0 {
				bufBefore.push(buf{line, lineNum})
			}

		}
	}

	if f.count {
		fmt.Fprintf(writer, "%d", count)
	}
}

func main() {

	f := &flags{}

	flag.IntVar(&f.after, "A", 0, "печатать +N строк после совпадения")
	flag.IntVar(&f.before, "B", 0, "печатать +N строк до совпадения")
	flag.IntVar(&f.context, "C", 0, "печатать ±N строк вокруг совпадения")
	flag.BoolVar(&f.count, "c", false, "количество строк")
	flag.BoolVar(&f.ignoreCase, "i", false, "игнорировать регистр")
	flag.BoolVar(&f.invert, "v", false, "вместо совпадения, исключать")
	flag.BoolVar(&f.fixed, "F", false, "точное совпадение со строкой")
	flag.BoolVar(&f.lineNum, "n", false, "печатать номер строки")

	flag.Parse()

	if flag.NArg() >= 2 {
		f.pattern = flag.Arg(0)

		r, err := os.Open(flag.Arg(1))
		if err != nil {
			panic(err)
		}

		w := os.Stdout
		if flag.NArg() > 2 {
			w, err = os.OpenFile(flag.Arg(2), os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				panic(err)
			}
		}
		grep(r, w, f)
	} else {
		fmt.Println("Usage: grep [-A <N>] [-B <N>] [-C <N>] [-c] [-i] [-v] [-F] [-n] <pattern> <file> [<file>]")
	}
}

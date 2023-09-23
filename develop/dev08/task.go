package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"

	gops "github.com/mitchellh/go-ps"
)

/*
=== Взаимодействие с ОС ===

Необходимо реализовать свой собственный UNIX-шелл-утилиту с поддержкой ряда простейших команд:

- cd <args> - смена директории (в качестве аргумента могут быть то-то и то)
- pwd - показать путь до текущего каталога
- echo <args> - вывод аргумента в STDOUT
- kill <args> - "убить" процесс, переданный в качесте аргумента (пример: такой-то пример)
- ps - выводит общую информацию по запущенным процессам в формате *такой-то формат*

Так же требуется поддерживать функционал fork/exec-команд

Дополнительно необходимо поддерживать конвейер на пайпах
(linux pipes, пример cmd1 | cmd2 | .... | cmdN).

Шелл — это обычная консольная программа, которая будучи
запущенной, в интерактивном сеансе выводит некое приглашение
в STDOUT и ожидает ввода пользователя через STDIN. Дождавшись
ввода, обрабатывает команду согласно своей логике
и при необходимости выводит результат на экран. Интерактивный
сеанс поддерживается до тех пор, пока не будет введена
команда выхода (например \quit).
*/

var (
	stderr = os.Stderr
	in     = os.Stdin
	out    = os.Stdout
)

func main() {
	shell()
}

func shell() {
	reader := bufio.NewReader(in)
	var w io.Writer = out

	for {
		head()

		byteLine, _, err := reader.ReadLine()
		if err != nil {
			_, _ = fmt.Fprint(stderr, err)
		}
		line := string(byteLine)

		commands := strings.Split(line, "|")

		buf := &strings.Builder{}

		if len(commands) > 1 {
			w = buf
		}

		for _, command := range commands {
			if command == "" {
				continue
			}

			cmd := strings.Fields(command)
			name := cmd[0]
			args := cmd[1:]

			if buf.Len() > 0 {
				args = append(args, buf.String())
				buf.Reset()
			}

			runCommand(w, name, args...)
		}

		w = out
		_, _ = fmt.Fprintln(w, buf.String())
	}
}

func head() {
	u, _ := user.Current()
	dir, _ := os.Getwd()
	prefix := ""
	if len(dir) > 20 {
		dir = dir[len(dir)-20:]
		prefix = "..."
	}
	fmt.Printf("%s: %s%s> ", u.Username, prefix, dir)
}

func runCommand(w io.Writer, name string, args ...string) {
	switch name {
	case "pwd":
		pwd(w)
	case "cd":
		cd(args...)
	case "echo":
		echo(w, args...)
	case "ps":
		ps(w)
	case "kill":
		kill(args...)
	case "exit", "quit":
		os.Exit(0)
	default:
		execCommand(w, name, args...)
	}
}

func cd(args ...string) {
	if len(args) != 1 {
		_, _ = fmt.Fprint(stderr, "cd <dir>")
		return
	}
	if err := os.Chdir(args[0]); err != nil {
		_, _ = fmt.Fprint(stderr, err)
	}
	return
}

func pwd(w io.Writer) {
	dir, err := os.Getwd()
	if err != nil {
		_, _ = fmt.Fprint(stderr, err)
		return
	}
	_, _ = fmt.Fprintf(w, "%s", dir)
}

func echo(w io.Writer, args ...string) {

	if len(args) == 0 {
		_, _ = fmt.Fprint(os.Stdout, "InputObject: ")
		stdin := bufio.NewReader(os.Stdin)
		obj, _, err := stdin.ReadLine()
		if err != nil {
			_, _ = fmt.Fprint(stderr, err)
			return
		}
		args = strings.Fields(string(obj))
	}

	for i := 0; i < len(args); i++ {
		if i != 0 {
			_, _ = fmt.Fprintf(w, "%s", " ")
		}

		res := args[i]
		if strings.HasPrefix(args[i], "$") {
			args[i] = strings.TrimPrefix(args[i], "$")
			res = os.Getenv(args[i])
		}
		_, _ = fmt.Fprintf(w, "%s", res)
	}
}

func ps(w io.Writer) {
	p, err := gops.Processes()
	if err != nil {
		_, _ = fmt.Fprint(stderr, err)
		return
	}
	_, _ = fmt.Fprintln(w, "\tPID\tPPID\tCMD")
	for _, process := range p {
		_, _ = fmt.Fprintf(w, "\t%v\t%v\t%v\n", process.Pid(), process.PPid(), process.Executable())
	}
}

func kill(args ...string) {
	if len(args) > 1 || len(args) == 0 {
		_, _ = fmt.Fprint(stderr, errors.New("kill<PID>"))
		return
	}

	pid, err := strconv.Atoi(args[0])
	if err != nil {
		_, _ = fmt.Fprint(stderr, err)
		return
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		_, _ = fmt.Fprint(stderr, err)
		return
	}

	err = process.Kill()
	if err != nil {
		_, _ = fmt.Fprint(stderr, err)
		return
	}
}

func execCommand(w io.Writer, name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = w
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		_, _ = fmt.Fprint(stderr, err)
		return
	}
}

package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

type config struct {
	timeout time.Duration
	host    string
	port    string
}

func main() {
	cfg := config{}

	flag.DurationVar(&cfg.timeout, "timeout", 10*time.Second, "timeout")

	flag.Parse()

	if len(flag.Args()) != 2 {
		log.Fatal("Usage: telnet [flags] host port")
	}

	cfg.host = flag.Args()[0]
	cfg.port = flag.Args()[1]

	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signals
		cancel()
	}()

	telnet(ctx, &cfg)
}

func telnet(ctx context.Context, cfg *config) {
	addr := fmt.Sprintf("%s:%s", cfg.host, cfg.port)
	conn, err := net.DialTimeout("tcp", addr, cfg.timeout)
	if err != nil {
		log.Fatalf("ошибка при подключении: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go readWrite(conn, os.Stdout, cancel)
	go readWrite(os.Stdin, conn, cancel)

	<-ctx.Done()
}

func readWrite(r io.Reader, w io.Writer, cancel context.CancelFunc) {
	rw := bufio.NewReadWriter(bufio.NewReader(r), bufio.NewWriter(w))
	defer rw.Flush()

	for {
		line, _, err := rw.ReadLine()
		if err != nil {
			if err == io.EOF {
				cancel()
				return
			}
			log.Fatalf("Ошибка при чтении: %v", err)
		}

		_, err = rw.Write(line)
		_, _ = rw.WriteRune('\n')
		if err != nil {
			log.Fatalf("Ошибка при записи: %v", err)
		}
		rw.Flush()
	}
}

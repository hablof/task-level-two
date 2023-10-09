package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	flag "github.com/spf13/pflag"
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

type telnetConfig struct {
	host    string
	port    string
	timeout time.Duration
}

var (
	ErrNotSetHostPort = errors.New("host/port unspecified")
)

func parseFlags(progname string, args []string) (cfg *telnetConfig, output string, err error) {
	flags := flag.NewFlagSet(progname, flag.ContinueOnError)
	buf := bytes.Buffer{}
	flags.SetOutput(&buf)

	to := flags.IntP("timeout", "t", 10, "sets TCP connection timeout in seconds")

	err = flags.Parse(args)
	if err != nil {
		return nil, buf.String(), err
	}

	host := strings.TrimSpace(flags.Arg(0))
	port := strings.TrimSpace(flags.Arg(1))

	if host == "" || port == "" {
		return nil, "", ErrNotSetHostPort
	}

	cfg = &telnetConfig{
		host:    host,
		port:    port,
		timeout: time.Second * time.Duration(*to),
	}

	return cfg, "", nil
}

func main() {

	cfg, output, err := parseFlags(os.Args[0], os.Args[1:])
	if err != nil {
		fmt.Println(output)
		fmt.Println(err)
		os.Exit(1)
	}

	conn, err := net.DialTimeout("tcp", cfg.host+":"+cfg.port, cfg.timeout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.Println("connected")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		defer cancel()
		defer fmt.Println("reader from conn exited")
		reader := bufio.NewReader(conn)
		for {
			s, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("closing conn\n")
				cancel()
				break
			}

			s = strings.TrimRight(s, "\r\n")

			fmt.Printf("-> %s\n", s)
		}
	}()

	stdInCh := make(chan []byte)
	go func() {
		defer cancel()
		defer log.Println("reader from stdin exited")

		reader := bufio.NewReader(os.Stdin)

		for {
			s, err := reader.ReadBytes('\n')
			if err != nil {
				log.Printf("err on read from stdin: %v\n", err)
				return
			}
			stdInCh <- s
		}

	}()

	// go func() {
	// 	scanner := bufio.NewScanner(os.Stdin)
	// 	for {
	// 		log.Println("looped")
	// 		scanner.Scan()
	// 		stdInCh <- []byte(scanner.Text())
	// 	}
	// }()

	go func() {
		defer cancel()
		defer log.Println("writer to conn exited")
		for {
			select {
			case s := <-stdInCh:
				if _, err := conn.Write(s); err != nil {
					log.Printf("err on write: %v\n", err)
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	terminationChannel := make(chan os.Signal, 1)

	// не знаю как на Windows обработать Ctrl+D, дай бог этот код работает на Linux/MacOS
	// при нажатии Ctrl+D в терминале просто появляется "^D"
	signal.Notify(terminationChannel, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-terminationChannel:
	case <-ctx.Done():
	}
	conn.Close()
	cancel()

	// не знаю почему, без Sleep терминал не хотел адекватно возвращать управление
	// приходилось нажимать Enter
	// использую git bash под windows
	time.Sleep(10 * time.Millisecond)
	// fmt.Println("Session ended. To exit send any message...")
	// os.Stdin.Close()
}

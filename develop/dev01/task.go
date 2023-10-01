package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

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

var (
	defaultTimeServer = "0.beevik-ntp.pool.ntp.org"
	helpMSG           = `-url <URL NTP-сервера> - задаёт `
)

func main() {
	var timeServerURL string

	flag.StringVar(&timeServerURL, "url", defaultTimeServer, "Определяет URL NTP-сервера")
	flag.ErrHelp = errors.New(helpMSG)

	flag.Parse()

	time, err := ntp.Time(timeServerURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось получить точное время от сервера %v: %v\n", timeServerURL, err)
		os.Exit(1)
	}

	fmt.Printf("Точное время: %v\n", time)
}

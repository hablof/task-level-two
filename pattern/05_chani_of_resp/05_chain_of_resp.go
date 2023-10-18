package main

import "fmt"

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/

type eventLevel uint8

const (
	_ eventLevel = iota
	trace
	regular
	warning
	fatal
)

type Logger interface {
	Log(mgs string)
}

type loggerChainElement struct {
	level   eventLevel
	current Logger
	next    *loggerChainElement
}

func (l loggerChainElement) handleLog(msg string, lvl eventLevel) {
	if lvl < l.level {
		return
	}

	l.current.Log(msg)

	if l.next != nil {
		l.next.handleLog(msg, lvl)
	}
}

// конкретная реализация логгера
type stdoutLogger struct{}

func (s stdoutLogger) Log(mgs string) {
	fmt.Printf("log to stdout: %s\n", mgs)
}

// конкретная реализация логгера
type emailNotifier struct {
	emails []string
}

func (e emailNotifier) Log(mgs string) {
	for _, email := range e.emails {
		fmt.Printf("send email to %s: %s\n", email, mgs)
	}
}

// конкретная реализация логгера

type ironLadyNightcall struct {
	phones []string
}

func (n ironLadyNightcall) Log(mgs string) {
	for _, phone := range n.phones {
		fmt.Printf("call to %s: %s\n", phone, mgs)
	}
}

func main() {
	ironLady := ironLadyNightcall{
		phones: []string{"+7 962 962 10 10"},
	}

	emails := emailNotifier{
		emails: []string{"hablof@yandex.ru"},
	}

	stdout := stdoutLogger{}

	loggerChain := loggerChainElement{
		level:   regular,
		current: stdout,
		next: &loggerChainElement{
			level:   warning,
			current: emails,
			next: &loggerChainElement{
				level:   0,
				current: ironLady,
				next:    nil,
			},
		},
	}

	loggerChain.handleLog("trace shouldn't be logged", trace)
	loggerChain.handleLog("regular should be logged once", regular)
	fmt.Println()
	loggerChain.handleLog("warning should be logged twice", warning)
	fmt.Println()
	loggerChain.handleLog("fatal should be logged thrice", fatal)
}

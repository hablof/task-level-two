package main

import "fmt"

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern
*/

type Strategy interface {
	Execute(a int, b int) int
}

// конкретная стратегия №1
type Multiplier struct{}

func (m Multiplier) Execute(a int, b int) int {
	return a * b
}

// конкретная стратегия №2
type Substractor struct{}

func (s Substractor) Execute(a int, b int) int {
	return a - b
}

// Структура использующая стратегии
type Context struct {
	s Strategy
}

func (c Context) Do(a, b int) {
	fmt.Println(c.s.Execute(a, b))
}

func main() {
	mul := Context{
		s: Multiplier{},
	}

	sub := Context{
		s: Substractor{},
	}

	mul.Do(10, 5)
	sub.Do(10, 5)
}

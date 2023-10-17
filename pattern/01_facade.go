package pattern

import "fmt"

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern
*/

// Проблема: мас нужна структура реализующая интерфейс
// с требуемыми методами
type demandedMethods interface {
	Method1() int
	Method2(a interface{})
	Method3(a []int) int
}

// При этом у нас уже имеются структуры Class1, Class2 и Class3
// которые вместе реализуют все необходимые методы интерфейса demandedMethods
type Class1 struct{}

func (c1 Class1) Method1() int {
	return 42
}

type Class2 struct{}

func (c2 Class2) Method2(a interface{}) {
	fmt.Println(a)
}

type Class3 struct{}

func (c3 Class3) Method3(a []int) int {
	sum := 0
	for _, elem := range a {
		sum += elem
	}

	return sum
}

// Реализуем структуру, готорая будет являться прослойкой, объединяющей
// имеющиеся структуры в одну -- фасад
type Facade struct {
	Class1
	Class2
	Class3
}

// структура Facade реализует интерфейс demandedMethods
var _ demandedMethods = Facade{}

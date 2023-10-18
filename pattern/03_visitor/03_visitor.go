package pattern

import "fmt"

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern
*/

/*
	Паттерн позволяет внедрить новую функциональность над существуемыми структурами, вынеся её реализацию в методы отжельной структуры

	Плюсы:
		для используемых структур не требует написания большого объёма кода
		S - принцип на внедряемую однородную функциональность
	Минусы:
		возможны проблемы с желанием использовать неэкспортируемые поля/методы структур
*/

// Имеются реализованные стркутуры MyStruct1 и MyStruct2
type MyStruct1 struct {
	a int
	b int
}

type MyStruct2 struct {
	a string
	b string
}

// Если предполагается внедрение только одного вида функциональности
// то интерфейс не обязателен
type Visitor interface {
	VisitMyStruct1(s MyStruct1)
	VisitMyStruct2(s MyStruct2)
}

// структура, которая будет предоставлять новую функциональность для MyStruct1 и MyStruct2
type Printer struct{}

func (Printer) VisitMyStruct1(s MyStruct1) {
	fmt.Println(s.a * s.b)
}

func (Printer) VisitMyStruct2(s MyStruct2) {
	fmt.Println(s.a + s.b)
}

// структуры MyStruct1 и MyStruct2 реализуют вспомогательный интерфейс accepter...
type accepter interface {
	accept(v Visitor)
}

func (m MyStruct1) accept(v Visitor) {
	v.VisitMyStruct1(m)
}

func (m MyStruct2) accept(v Visitor) {
	v.VisitMyStruct2(m)
}

func Example() {
	// блягодаря этому мы можем пробежаться по ним как-будто они имеют общую функциональность
	printing := []accepter{MyStruct1{}, MyStruct2{}}

	for _, elem := range printing {
		elem.accept(Printer{})
	}
}

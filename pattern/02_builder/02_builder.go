package pattern

/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern
*/

/*
	Паттерн позволяет создавать объекты пошагово и ветвить задаваемые параметры в зависимости от внешних условий

	Пример использования на практике: для динамического посторения SQL-запростов я использовал библиотеку https://github.com/Masterminds/squirrel
	с ней очень легко подготовить запрос для множественного INSERT'а, например
*/

type enum uint8

const (
	_ enum = iota
	firstOption
	secondOption
)

// сложный объект
type product struct {
	strField string
	params   []string
	intVal   int
	oneOf    enum
}

// строитель
type builder struct {
	product product
}

func (b *builder) newProduct() *builder {
	b.product = product{}
	return b
}

func (b *builder) setStrField(str string) *builder {
	b.product.strField = str
	return b
}

func (b *builder) addParam(param string) *builder {
	b.product.params = append(b.product.params, param)
	return b
}

func (b *builder) setIntValue(val int) *builder {
	b.product.intVal = val
	return b
}

func (b *builder) setEnum(e enum) *builder {
	b.product.oneOf = e
	return b
}

// возвращает объект
func (b *builder) result() product {
	return b.product
}

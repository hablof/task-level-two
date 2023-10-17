Что выведет программа? Объяснить вывод программы. Объяснить внутреннее устройство интерфейсов и их отличие от пустых интерфейсов.

```go
package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)
	fmt.Println(err == nil)
}
```

Ответ:
```
<nil>
false

```
Потому что интерфейс -- структура из двух полей: тип, ссылка на структуру. Интерфейс равен nil если оба этих поля нулевые.
Сторочкой `var err *os.PathError = nil` было задано значение полю типа.
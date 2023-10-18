package pattern

/*
	Реализовать паттерн «комманда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern
*/

/*
	Паттерн представляет необходимые действия в виде объекта, изменяющего состояния другого объекта

	Плюсы:
		1) Поскольку оъект хранит всю информацию о необходимых действиях, исполнение команды можно отложить
		2) Позволяет реализовать отмену операций (как в моём примере)
	Минусы:
		1) Усложнение
*/

type Command interface {
	Execute()
	Unexecute()
}

// структура к которой применяются комманды
type Receiver struct {
	fieldToChange int
}

// команда для инкрементирования поля receiver'а
type incrementReciever struct {
	obj *Receiver
}

// конструктор
func NewInc(obj *Receiver) incrementReciever {
	return incrementReciever{
		obj: obj,
	}
}

func (i incrementReciever) Execute() {
	// логика домена
	i.obj.fieldToChange++
}

func (i incrementReciever) Unexecute() {
	// логика домена
	i.obj.fieldToChange--
}

// команда для удвоения поля receiver'а
type doubleReciever struct {
	obj *Receiver
}

// конструктор
func NewDouble(obj *Receiver) doubleReciever {
	return doubleReciever{
		obj: obj,
	}
}

func (d doubleReciever) Execute() {
	// логика домена
	d.obj.fieldToChange *= 2
}

func (d doubleReciever) Unexecute() {
	// логика домена
	d.obj.fieldToChange /= 2
}

type Invoker struct {
	// стек команд
	cmds []Command
}

func (i *Invoker) Execute(c Command) {
	c.Execute()
	i.cmds = append(i.cmds, c)
}

func (i *Invoker) Undo() {
	if len(i.cmds) == 0 {
		return
	}

	i.cmds[len(i.cmds)-1].Unexecute()
	i.cmds = i.cmds[:len(i.cmds)-1]
}

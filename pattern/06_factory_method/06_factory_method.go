package pattern

import "errors"

/*
	Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern
*/

/*
	позволяет переиспользовать код создания сложного объекта
*/

type seat struct {
	hasBack      bool
	hasArmrests  bool
	weigth       int
	height       int
	productivity int
	comfort      int
}

type seatType uint8

const (
	_ seatType = iota
	stool
	chair
	armchair
	computerChair
)

type seatFactory struct{}

func (sf seatFactory) makeSeat(known seatType) (seat, error) {
	switch known {
	case stool:
		return seat{
			hasBack:      false,
			hasArmrests:  false,
			weigth:       2,
			height:       45,
			productivity: 8,
			comfort:      2,
		}, nil

	case chair:
		return seat{
			hasBack:      true,
			hasArmrests:  false,
			weigth:       3,
			height:       50,
			productivity: 10,
			comfort:      5,
		}, nil

	case armchair:
		return seat{
			hasBack:      true,
			hasArmrests:  true,
			weigth:       20,
			height:       40,
			productivity: 2,
			comfort:      10,
		}, nil

	case computerChair:
		return seat{
			hasBack:      true,
			hasArmrests:  true,
			weigth:       8,
			height:       60,
			productivity: 8,
			comfort:      5,
		}, nil
	}

	return seat{}, errors.New("unknown seat")
}

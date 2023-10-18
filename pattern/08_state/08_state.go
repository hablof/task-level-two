package main

import (
	"fmt"
	"math/rand"
)

/*
	Реализовать паттерн «состояние».

Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.

	https://en.wikipedia.org/wiki/State_pattern
*/

/*
	паттерн разделяет поведение объекта в разных состояниях на различные классы
	Плюсы:
		1) инкапсуляция состояния
	Минусы:
		2) нужен мета-объект управляющий изменением состояний
*/

const threshold_CV = 70

// вспомогательный объект
type Battery struct {
	chargeLevel int
}

// интерфейс который должны выполнять все состояния
type charger interface {
	charge(b *Battery) charger
}

// зарядник, меняющий режимы своей работы
type ChargerContext struct {
	state charger
}

func (c *ChargerContext) Charge(b *Battery) {
	if newState := c.state.charge(b); newState != nil {
		c.state = newState
	}
}

func (c *ChargerContext) IsCharged(b *Battery) bool {
	if b.chargeLevel >= 100 {
		return true
	}

	return false
}

// состояние зарядки постоянным током
type CCmode struct{}

func (cc CCmode) charge(b *Battery) charger {
	if b.chargeLevel > threshold_CV {
		return CVmode{}
	}

	if rand.Intn(100) < 10 {
		fmt.Println("charger overheated!")
		return &OHmode{}
	}

	fmt.Printf("battery charge level %d%%, charging CC...\n", b.chargeLevel)
	b.chargeLevel += 10
	return nil
}

// состояние зарядки постоянным напряжением
type CVmode struct{}

func (cv CVmode) charge(b *Battery) charger {
	if b.chargeLevel <= threshold_CV {
		return CCmode{}
	}

	if rand.Intn(100) < 2 {
		fmt.Println("charger overheated!")
		return &OHmode{}
	}

	fmt.Printf("battery charge level %d%%, charging CV...\n", b.chargeLevel)
	if b.chargeLevel < 90 {
		b.chargeLevel += 7
	} else if b.chargeLevel < 97 {
		b.chargeLevel += 3
	} else {
		b.chargeLevel += 1
	}
	return nil
}

// защитное состояние
type OHmode struct {
	coolingStage int
}

func (oh *OHmode) charge(b *Battery) charger {
	if oh.coolingStage > 2 {
		return CCmode{}
	}
	fmt.Println("overheat protection, please wait...")
	oh.coolingStage++
	return nil
}

func main() {
	b := Battery{
		chargeLevel: 23,
	}

	charger := ChargerContext{
		state: CCmode{},
	}

	for !charger.IsCharged(&b) {
		charger.Charge(&b)
	}

	fmt.Println("battery charged!")
}

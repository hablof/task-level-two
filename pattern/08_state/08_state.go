package main

import (
	"errors"
	"fmt"
	"math/rand"
)

/*
	Реализовать паттерн «состояние».

Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.

	https://en.wikipedia.org/wiki/State_pattern
*/
const threshold_CV = 70

var ErrOverheat = errors.New("charger overheated")

type Battery struct {
	chargeLevel int
}

type Charger interface {
	Charge(b *Battery) Charger
}

type ChargerContext struct {
	state Charger
}

func (c *ChargerContext) Charge(b *Battery) {
	if newState := c.state.Charge(b); newState != nil {
		c.state = newState
	}
}

func (c *ChargerContext) IsCharged(b *Battery) bool {
	if b.chargeLevel >= 100 {
		return true
	}

	return false
}

type CCmode struct{}

func (cc CCmode) Charge(b *Battery) Charger {
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

type CVmode struct{}

func (cv CVmode) Charge(b *Battery) Charger {
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

type OHmode struct {
	coolingStage int
}

func (oh *OHmode) Charge(b *Battery) Charger {
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

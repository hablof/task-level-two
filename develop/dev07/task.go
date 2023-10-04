package main

import (
	"fmt"
	"sync"
	"time"
)

/*
=== Or channel ===

Реализовать функцию, которая будет объединять один или более done каналов в single канал если один из его составляющих каналов закроется.
Одним из вариантов было бы очевидно написать выражение при помощи select, которое бы реализовывало эту связь,
однако иногда неизестно общее число done каналов, с которыми вы работаете в рантайме.
В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or каналов, реализовывала весь функционал.
*/

var or func(channels ...<-chan interface{}) <-chan interface{}

// Сначала я написал такую функцию, но потом понял,
// что в таком использовании есть шанс нарваться на панику
// двойного закрытия каналов, если "одновременно" будут закрыты два канала и числа переданных в функцию
// в горутинах селект провалится в ветку close(cancelChannel), close(outChannel)
func BAD_FUNC(channels ...<-chan interface{}) <-chan interface{} {
	cancelChannel := make(chan interface{})

	outChannel := make(chan interface{})

	if len(channels) == 0 {
		defer close(outChannel)
	}

	for _, ch := range channels {
		go func(ch <-chan interface{}) {
			for {
				select {
				case <-cancelChannel:
					return

				case _, ok := <-ch:
					if !ok {
						close(cancelChannel)
						close(outChannel)

						return
					}
				}
			}
		}(ch)
	}

	return outChannel
}

// "объединять один или более DONE каналов" - значит передача данных по каналу не предполагается,
// а каналы будут использоваться исключительно в качестве способа сообщить Done.
// Тем не менее, перестрахуемся и при чтении _, ok := <- ch, ok == true не будем закрывать объединяющий канал
//
// В этом варианте соблюдается следующий принцип взаимоотношения горутин и каналов:
// горутина ЛИБО читает из канала/обрабатывает закрытие, ЛИБО пишет в канал/закрывает его
func orChannel(channels ...<-chan interface{}) <-chan interface{} {
	sharedCancelChannel := make(chan interface{})

	outChannel := make(chan interface{})

	wg := &sync.WaitGroup{}

	unaryCancelChannelsChunck := make([]chan interface{}, 0, len(channels))

	for _, ch := range channels {
		wg.Add(1)

		unaryCancelChannel := make(chan interface{})
		unaryCancelChannelsChunck = append(unaryCancelChannelsChunck, unaryCancelChannel)

		go func(inputCh <-chan interface{}) {
			defer wg.Done()

			for {
				select {
				case _, ok := <-inputCh:
					if !ok {
						select {
						case <-sharedCancelChannel:
							return

						case <-unaryCancelChannel:
							return
						}
					}

				case <-unaryCancelChannel:
					return
				}
			}
		}(ch)
	}

	go func() {
		sharedCancelChannel <- nil
		for _, ch := range unaryCancelChannelsChunck {
			close(ch)
		}
		wg.Wait()
		close(outChannel)
	}()

	if len(channels) == 0 {
		<-sharedCancelChannel
	}

	return outChannel
}

func main() {

	or = orChannel

	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("fone after %v", time.Since(start))
}

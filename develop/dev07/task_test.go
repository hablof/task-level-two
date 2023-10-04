package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_orChannel(t *testing.T) {

	t.Run("check time between", func(t *testing.T) {

		firstCh := make(chan interface{})
		secondCh := make(chan interface{})
		gotCh := orChannel(firstCh, secondCh)

		start := time.Now()
		go func() {
			<-time.After(100 * time.Millisecond)
			close(firstCh)
		}()

		go func() {
			<-time.After(time.Second)
			close(secondCh)
		}()

		<-gotCh
		dur := time.Since(start)

		assert.Less(t, dur, time.Second)
		assert.Greater(t, dur, 100*time.Millisecond)
	})

	t.Run("try catch panic on concurrent input channels close", func(t *testing.T) {
		defer func() {
			recovery := recover()
			assert.Nil(t, recovery)
		}()

		deadline := time.Now().Add(500 * time.Millisecond)
		ctxbg := context.Background()

		channels := make([]<-chan interface{}, 0, 500)
		for i := 0; i < 500; i++ {
			ctx, cf := context.WithDeadline(ctxbg, deadline)
			defer cf()

			inerfaceCh := make(chan interface{})
			go func() {
				<-ctx.Done()
				close(inerfaceCh)
			}()

			channels = append(channels, inerfaceCh)
		}

		<-orChannel(channels...)
	})
}

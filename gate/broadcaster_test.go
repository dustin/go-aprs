package main

import (
	"sync"
	"testing"

	"github.com/dustin/go-aprs"
)

func TestBroadcast(t *testing.T) {
	ch := make(chan aprs.APRSData)
	wg := sync.WaitGroup{}

	b := NewBroadcaster(ch)

	for i := 0; i < 5; i++ {
		wg.Add(1)

		cch := make(chan aprs.APRSData)

		b.Register(cch)

		go func() {
			defer wg.Done()
			defer b.Unregister(cch)
			<-cch
		}()

	}

	b.broadcast(aprs.APRSData{})

	wg.Wait()
}

package main

import (
	"github.com/dustin/go-aprs"
)

type broadcaster struct {
	input <-chan aprs.APRSData
	reg   chan chan<- aprs.APRSData
	unreg chan chan<- aprs.APRSData

	outputs map[chan<- aprs.APRSData]bool
}

func (b *broadcaster) cleanup() {
	for ch := range b.outputs {
		close(ch)
	}
}

func (b *broadcaster) broadcast(m aprs.APRSData) {
	for ch := range b.outputs {
		ch <- m
	}
}

func (b *broadcaster) run() {
	defer b.cleanup()

	for {
		select {
		case m, ok := (<-b.input):
			if ok {
				b.broadcast(m)
			} else {
				return
			}
		case ch, ok := (<-b.reg):
			if ok {
				b.outputs[ch] = true
			} else {
				return
			}
		case ch := (<-b.unreg):
			delete(b.outputs, ch)
		}
	}
}

func NewBroadcaster(input <-chan aprs.APRSData) *broadcaster {
	b := &broadcaster{
		input:   input,
		reg:     make(chan chan<- aprs.APRSData),
		unreg:   make(chan chan<- aprs.APRSData),
		outputs: make(map[chan<- aprs.APRSData]bool),
	}

	go b.run()

	return b
}

func (b *broadcaster) Register(newch chan<- aprs.APRSData) {
	b.reg <- newch
}

func (b *broadcaster) Unregister(newch chan<- aprs.APRSData) {
	b.unreg <- newch
}

func (b *broadcaster) Close() error {
	close(b.reg)
	return nil
}

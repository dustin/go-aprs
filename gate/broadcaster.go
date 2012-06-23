package main

import (
	"github.com/dustin/go-aprs"
)

type broadcaster struct {
	input <-chan aprs.APRSMessage
	reg   chan chan<- aprs.APRSMessage
	unreg chan chan<- aprs.APRSMessage

	outputs map[chan<- aprs.APRSMessage]bool
}

func (b *broadcaster) cleanup() {
	for ch := range b.outputs {
		close(ch)
	}
}

func (b *broadcaster) broadcast(m aprs.APRSMessage) {
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

func NewBroadcaster(input <-chan aprs.APRSMessage) *broadcaster {
	b := &broadcaster{
		input:   input,
		reg:     make(chan chan<- aprs.APRSMessage),
		outputs: make(map[chan<- aprs.APRSMessage]bool),
	}

	go b.run()

	return b
}

func (b *broadcaster) Register(newch chan<- aprs.APRSMessage) {
	b.reg <- newch
}

func (b *broadcaster) Unregister(newch chan<- aprs.APRSMessage) {
	b.unreg <- newch
}

func (b *broadcaster) Close() error {
	close(b.reg)
	return nil
}

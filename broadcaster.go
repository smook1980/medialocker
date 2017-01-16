package medialocker

import (
	"sync"
	"errors"
	"time"
)

type Broadcaster struct {
	lock sync.RWMutex
	input chan interface{}
	stopped chan interface{}
	listeners []*listener
}

func NewBroadcaster(inBuffSize uint) *Broadcaster {
	bc := Broadcaster{input: make(chan interface{}, inBuffSize), stopped: make(chan interface{}, 1)}

	go bc.relay()

	return &bc
}

func (bc *Broadcaster) IsDestroyed() bool {
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	return bc.listeners == nil
}

func (bc *Broadcaster) Destroy(timeout int) error {
	bc.lock.Lock()
	defer bc.lock.Unlock()
	if bc.input == nil {
		return errors.New("Calling BroadcastChan.Destroy multiple times on an instance.")
	}

	close(bc.input)

	var err error

	select {
	case <-bc.stopped:
	case <-time.After(time.Second * time.Duration(timeout)):
	  err = errors.New("BroadcastChan.Destroy timed out, forcing closing channels!")
	}

	for _, l := range bc.listeners {
		l.Close()
	}

	bc.listeners = nil

	return err
}

func (bc *Broadcaster) relay() {
	for msg := range bc.input {
		bc.broadcast(msg)
	}

	bc.input = nil
	bc.stopped<- struct{}{}
	close(bc.stopped)
}

func (bc *Broadcaster) sweep() {
	var listeners []*listener

	bc.lock.Lock()
	defer bc.lock.Unlock()

	for _, l := range(bc.listeners) {
		if !l.IsClosed() {
			listeners = append(listeners, l)
		}
	}

	bc.listeners = listeners
}

func (bc *Broadcaster) Listen(buffSize uint) (<-chan interface{}, Closer) {
	l := &listener{channel: make(chan interface{}, buffSize), closed: false}

	bc.lock.Lock()
	defer bc.lock.Unlock()

	closer := func() error {
		err := l.Close()
		bc.sweep()

		return err
	}

	bc.listeners = append(bc.listeners, l)

	return l.channel, closer
}

func (bc *Broadcaster) Send(msg interface{}) {
	bc.input<- msg
}

func (bc *Broadcaster) SendChannel() chan<- interface{} {
	return bc.input
}

func (bc *Broadcaster) broadcast(msg interface{}) {
	bc.lock.RLock()
	defer bc.lock.RUnlock()

	for _, ch := range bc.listeners {
		ch.Write(msg)
	}
}

type listener struct {
	lock sync.RWMutex
	closed bool
	channel chan interface{}
}

func (l *listener) Write(msg interface{}) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	if !l.closed {
		l.channel<- msg
	}
}

func (l *listener) IsClosed() bool {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.closed
}

func (l *listener) Close() error {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.closed {
		return errors.New("Called close on an already closed Broadcast listener.")
	}

	close(l.channel)
	l.closed = true

	return nil
}

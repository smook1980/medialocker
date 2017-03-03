package medialocker

import (
	"fmt"
	"testing"
	"time"
)

func TestBroadcaster(t *testing.T) {
	subject := NewBroadcaster(1)
	listen, _ := subject.Listen(0)

	msg := "This is a message!"
	subject.Send(msg)

	select {
	case output := <-listen:
		if output != msg {
			t.Errorf("Expected input message %s to equal output %s", msg, output)
		}
	case <-time.After(time.Second * time.Duration(1)):
		t.Error("Timed out waiting for msg from broadcaster!")
	}

	subject.Destroy(1)
	if !subject.IsDestroyed() {
		t.Error("IsDestroyed returned false after call to Destory, expected true!")
	}
}

func TestBroadcaster_CloserFunc(t *testing.T) {
	subject := NewBroadcaster(0)
	defer subject.Destroy(0)

	_, closer := subject.Listen(0)

	if err := closer(); err != nil {
		t.Errorf("Call to closer func returned err: %s", err)
	}

	if len(subject.listeners) != 0 {
		t.Errorf("Expected closer func to remove listener from broadcaster.")
	}

	_, closer = subject.Listen(1)

	subject.Send("Closer Called With Message Queued")

	if err := closer(); err != nil {
		t.Errorf("Call to closer func returned err: %s", err)
	}

	if len(subject.listeners) != 0 {
		t.Errorf("Expected closer func to remove listener from broadcaster.")
	}

	_, closer1 := subject.Listen(1)
	_, closer2 := subject.Listen(1)

	subject.Send("Closer does not remove active channels.")

	if err := closer1(); err != nil {
		t.Errorf("Call to closer func returned err: %s", err)
	}

	if count := len(subject.listeners); count != 1 {
		t.Errorf("Given two listeners, after closing one expected one to remain, but got: %v", count)
	}

	if err := closer2(); err != nil {
		t.Errorf("Call to closer func returned err: %s", err)
	}

	if count := len(subject.listeners); count != 0 {
		t.Errorf("Given two listeners, after closing two expected none to remain, but got: %v", count)
	}
}

func TestBroadcaster_Send(t *testing.T) {
	listenerCount := 10
	numMsgs := uint(5)

	var listeners []<-chan interface{}
	subject := NewBroadcaster(0)
	defer subject.Destroy(0)

	// Send works with no listeners...
	subject.Send("Message")

	for x := 0; x < listenerCount; x++ {
		it, _ := subject.Listen(numMsgs)
		listeners = append(listeners, it)
	}

	for x := uint(0); x < numMsgs; x++ {
		subject.Send(fmt.Sprintln("Message ", x))
	}

	for lx, l := range listeners {
		for x := uint(0); x < numMsgs; x++ {
			select {
			case msg := <-l:
				expectedMsg := fmt.Sprintln("Message ", x)
				if msg != expectedMsg {
					t.Errorf("Listner %v: Expected recieved message %s to equal %s", lx, msg, expectedMsg)
				}
			case <-time.After(time.Second * time.Duration(1)):
				t.Errorf("Listner %v: Timed out waiting for message number %v.", lx, x)
			}
		}
	}
}

package pubsub

import (
	"fmt"
	"testing"
	"time"
)

func Proc(i int) func(e *Event) {
	return func(e *Event) {
		time.Sleep(1 * time.Second)
		fmt.Println(i, e.Objects)
	}
}

func TestDirsFromPages(t *testing.T) {
	go Hub.Start(2)
	Hub.Sub("evt1", Proc(1), Proc(2), Proc(3))
	go Hub.Pub("evt1", 1)
	go Hub.Pub("evt1", 2)
	go Hub.Pub("evt1", 3)
	go Hub.Pub("evt1", 4)
	time.Sleep(10 * time.Second)
}

// simple publish <-> subscribe events processing for go
//
// Usage example
//	package main
//
//	import (
//		"fmt"
//		"time"
//
//		"github.com/vtg/pubsub"
//	)
//
//	var events = *pubsub.Hub
//
//	func main() {
//		events.Start(1)
//
//		events.Sub("CalcComplete", notifyComplete, notifyCompleteTime)
//		go runCalculation(2)
//		go runCalculation(4)
//		time.Sleep(10 * time.Second)
//	}
//
//	func runCalculation(n int) {
//		//processing some long calculation here
//		time.Sleep(2 * time.Second)
//		var r int
//		r = n * n
//		events.Pub("CalcComplete", r, time.Now())
//	}
//
//	func notifyComplete(e *pubsub.Event) {
//		fmt.Println("result: ", e.Objects[0])
//	}
//
//	func notifyCompleteTime(e *pubsub.Event) {
//		fmt.Println("result calculated at: ", e.Objects[1])
//	}
//
package pubsub

// Event structure to store event information
type Event struct {
	Name    string
	Objects []interface{}
}

type listener struct {
	event string
	funcs []func(*Event)
}

type HUB struct {
	pub       chan Event
	sub       chan listener
	listeners []listener
}

// Hub variable to share pubsub events
var Hub = &HUB{
	pub:       make(chan Event, 5),
	sub:       make(chan listener, 5),
	listeners: []listener{},
}

// Pub publish new event
func (h *HUB) Pub(name string, objects ...interface{}) {
	h.pub <- Event{name, objects}
}

// Sub subscribe for event
func (h *HUB) Sub(name string, f ...func(*Event)) {
	h.sub <- listener{name, f}
}

// Start pubsub process with max number of parallel processing
func (h *HUB) Start(num int) {
	if num == 0 {
		num = 1
	}
	go h.runSub()
	for i := 0; i < num; i++ {
		go h.runPub()
	}
}

func (h *HUB) runPub() {
	for {
		select {
		case c := <-h.pub:
			for _, v := range h.listeners {
				if v.event == c.Name {
					for _, f := range v.funcs {
						f(&c)
					}
				}
			}
		}
	}
}

func (h *HUB) runSub() {
	for {
		select {
		case s := <-h.sub:
			h.listeners = append(h.listeners, s)
		}
	}
}

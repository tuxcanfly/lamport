package main

import (
	"log"
	"math/rand"
	"sync/atomic"
	"time"
)

// Event is the message send event between two processes
type Event struct {
	from    *Process
	to      *Process
	message string
	time    uint32
}

// Process represents a distributed node that can send
// and receive messages between other nodes
type Process struct {
	name   string
	time   uint32
	events chan *Event
}

// Print the process name on formatting
func (p *Process) String() string {
	return p.name
}

// step increments the logical timestamp of the process
// using an atomic counter to ensure thread safety
func (p *Process) step() {
	atomic.AddUint32(&p.time, 1)
}

// send a message to the given process
func (p *Process) send(to *Process, message string) {
	p.step()
	e := &Event{
		from:    p,
		to:      to,
		message: message,
		time:    p.time,
	}
	to.events <- e
	log.Printf("SEND: %v => %v: %v", e.from, e.to, e.message)
}

// receive listens to messages sent to this process
// it runs as a goroutine
func (p *Process) receive() {
	for {
		select {
		case e := <-p.events:
			if p.time < e.time {
				p.time = e.time
			}
			p.step()
			e.time = p.time
		default:
		}
	}
}

// NewProcess returns a new process
// it triggers the receive goroutine
func NewProcess(name string) *Process {
	p := &Process{
		name:   name,
		events: make(chan *Event),
	}
	go p.receive()
	return p
}

func main() {

	// create three processes
	p1 := NewProcess("Alice")
	p2 := NewProcess("Bob")
	p3 := NewProcess("Charlie")

	processes := [3]*Process{p1, p2, p3}
	for {
		time.Sleep(time.Second)
		go func() {
			// pick 2 at random and send a message
			a := processes[rand.Int()%len(processes)]
			b := processes[rand.Int()%len(processes)]
			a.send(b, "hi")
			log.Printf("%v[%v], %v[%v], %v[%v]", p1, p1.time, p2, p2.time, p3, p3.time)
		}()
	}
}

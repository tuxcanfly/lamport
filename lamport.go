package main

import (
	"log"
	"math/rand"
	"sync/atomic"
	"time"
)

type Event struct {
	from    *Process
	to      *Process
	message string
	time    uint32
}

type Process struct {
	name   string
	time   uint32
	events chan *Event
}

func (p *Process) String() string {
	return p.name
}

func (p *Process) step() {
	atomic.AddUint32(&p.time, 1)
}

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

func NewProcess(name string) *Process {
	p := &Process{
		name:   name,
		events: make(chan *Event),
	}
	go p.receive()
	return p
}

func main() {
	p1 := NewProcess("Alice")
	p2 := NewProcess("Bob")
	p3 := NewProcess("Charlie")

	processes := [3]*Process{p1, p2, p3}
	for {
		time.Sleep(time.Second)
		go func() {
			a := processes[rand.Int()%len(processes)]
			b := processes[rand.Int()%len(processes)]
			a.send(b, "hi")
			log.Printf("%v[%v], %v[%v], %v[%v]", p1, p1.time, p2, p2.time, p3, p3.time)
		}()
	}
}

package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

var (
	server string = "10.0.0.210:8080"
)

// Boa Constrictor, a poem by Shel Silverstein
var poem = `Oh, I'm being eaten
By a boa constrictor,
A boa constrictor,
A boa constrictor,
I'm being eaten by a boa constrictor,
And I don't like it--one bit.
Well, what do you know?
It's nibblin' my toe.
Oh, gee,
It's up to my knee.
Oh my,
It's up to my thigh.
Oh, fiddle,
It's up to my middle.
Oh, heck,
It's up to my neck.
Oh, dread,
It's upmmmmmmmmmmffffffffff . . .`

func segment(in chan []byte, out chan []byte) {
	var buf [512]byte
	for {
		c, err := net.Dial("tcp", server)
		for ; err != nil; c, err = net.Dial("tcp", server) {
			println(err.Error())
			time.Sleep(5 * time.Second)
		}
		for {
			select {
			case msg := <-in:
				_, err := c.Write(msg)
				if err != nil {
					log.Fatal(err.Error())
				}
				time.Sleep(100 * time.Millisecond)
				n, err := c.Read(buf[:])
				if err != nil {
					log.Fatal(err.Error())
				}
				out <- buf[:n]
			}
		}
	}
}

func feedit(head chan []byte) {
	for i := 0; i < 100; i++ {
		head <- []byte(fmt.Sprintf("\n---%d---", i))
		for _, line := range strings.Split(poem, "\n") {
			head <- []byte(line)
		}
	}
}

var head = make(chan []byte)
var a = make(chan []byte)
var b = make(chan []byte)
var c = make(chan []byte)
var d = make(chan []byte)
var e = make(chan []byte)
var f = make(chan []byte)
var tail = make(chan []byte)

func main() {

	// The snake
	go segment(head, a)
	go segment(a, b)
	go segment(b, c)
	go segment(c, d)
	go segment(d, e)
	go segment(e, f)
	go segment(f, tail)

	go feedit(head)

	for {
		select {
		case msg := <-tail:
			println(string(msg))
		}
	}
}

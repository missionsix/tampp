/* Editor Settings: expandtabs and use 4 spaces for indentation
 * ex: set softtabstop=4 tabstop=8 expandtab shiftwidth=4: *
 * -*- mode: c, c-basic-offset: 4 -*- */

/* Copyright (C) Patrick Andrew. All Rights Reserved
 *
 * Authors: Patrick Andrew <missionsix@gmail.com>
 */
package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"semaphore"
)

type Chopstick struct {	
	sem semaphore.Semaphore
	owner *int
}

func NewChopstick() *Chopstick {
	return &Chopstick{ sem: make(semaphore.Semaphore, 1), owner: nil }
}

func (c *Chopstick) Signal(n *int) {
	c.owner = nil
	c.sem.V(1)
}

func (c *Chopstick) Wait(n *int) {
	c.sem.P(1)
	c.owner = n
}

func main() {

	// grab Args without prog name
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Usage: ./%s <num resources>", os.Args[0])
		os.Exit(1)
	}

	num_philosophers, err := strconv.Atoi(args[0])
	if (err != nil) {
		fmt.Println("Usage: ./%s <num resources>", os.Args[0])
		os.Exit(1)
	}
	sigs := make(chan os.Signal, 1)
	done := make(chan bool)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Starting...")

	// simple signal handler goroutine
	go func() {
		fmt.Println("Sig handler launched")
		sig := <- sigs //blocks on sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true		
	}()

	n_chops := (num_philosophers)
	chopsticks := make( []*Chopstick, n_chops)
	for i := 0; i < n_chops; i++ {
		chopsticks[i] = NewChopstick()
	}

	chairs := make( semaphore.Semaphore, n_chops - 1 )

	// Philosopher go routine
	for i := 0; i < num_philosophers; i++ {
		var n int = i
		var first int
		var second int
		go func() {
			for ;; {
				select {
				case complete := <- done:
					if complete {
						break
					}
				default:
				}

				// get chopsticks
				first = n
				second = (n + 1) % n_chops

				// try to grab a seat
				chairs.P(1)

				chopsticks[first].Wait(&n)
				chopsticks[second].Wait(&n)

				fmt.Printf("Philosopher %d eating...\n", n)
				for j := range chopsticks {
					fmt.Printf("\t%d: %#p", j, chopsticks[j].owner)
				}
				fmt.Println()

				time.Sleep(time.Duration(rand.Intn(num_philosophers)) * 100)
				chopsticks[second].Signal(&n)
				chopsticks[first].Signal(&n)

				// stand up
				chairs.V(1)

				//think
				time.Sleep(time.Duration(rand.Intn(num_philosophers)) * 100)
			}
		}()
	}

	// wait for SIGINT
	<- done
}

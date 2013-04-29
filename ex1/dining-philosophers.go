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
	"sync/atomic"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unsafe"
)

type Chopstick struct {	
	owner *int
	owned bool
}

func NewChopstick() *Chopstick {
	return &Chopstick{ owner: nil, owned: false }
}

func (c *Chopstick) Take(n *int) bool {
	udest := (*unsafe.Pointer)(unsafe.Pointer(&c.owner))
	return atomic.CompareAndSwapPointer(udest,
		unsafe.Pointer(nil),
		unsafe.Pointer(n));
}

func (c *Chopstick) Return(n *int) bool {
	udest := (*unsafe.Pointer)(unsafe.Pointer(&c.owner))
	return atomic.CompareAndSwapPointer(udest,
		unsafe.Pointer(n),
		unsafe.Pointer(nil));
}

func main() {

	sigs := make(chan os.Signal, 1)
	done := make(chan bool)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	num_philosophers  := 5
	var i int

	fmt.Println("Starting...")

	// simple signal handler goroutine
	go func() {
		fmt.Println("Sig handler launched")
		sig := <- sigs //blocks on sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
		
	}()

	logger := make(chan []byte)
	go func() {
		for {
			select {
			case complete := <- done:
				if complete {
					break
				}
			default:
			}

			b := <- logger

			if len(b) == 0 {
				break
			}

			fmt.Println(b)
		}
	}()

	// get number of philosophers @ table
	// _, err = fmt.Scanf("%d"
	n_chops := (num_philosophers)
	chopsticks := make([]Chopstick, n_chops)

	// Philosopher go routine
	for i = 0; i < num_philosophers; i++ {
		var n int = i
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
				leftOf := n
				rightOf := (n + 1) % n_chops
				if (chopsticks[leftOf].Take(&n)) {
					if (chopsticks[rightOf].Take(&n)) {
						fmt.Printf("Philosopher %d eating...\n", n)
						for j := range chopsticks {
							fmt.Printf("\t%d: %#p", j, chopsticks[j].owner)
						}
						fmt.Println()
						time.Sleep(time.Duration(rand.Intn(num_philosophers)) * time.Second)
						chopsticks[rightOf].Return(&n)
					}
					chopsticks[leftOf].Return(&n)
				}
				//think
				time.Sleep(time.Duration(rand.Intn(num_philosophers)) * time.Second)
			}
		}()
	}

	<- done

	//mutex.Lock()
	
}




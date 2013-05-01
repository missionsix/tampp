/* Editor Settings: expandtabs and use 4 spaces for indentation
 * ex: set softtabstop=4 tabstop=8 expandtab shiftwidth=4: *
 * -*- mode: c, c-basic-offset: 4 -*- */

/*
 * Authors: Patrick Andrew <missionsix@gmail.com>
 */

// Package Semaphore provides a high level construct of semaphores
// implemented using Go channels.
package semaphore

// Empty is an empty struct acting as a single resource
type empty struct {}

// Semaphore is the actual type. Declare with := make(Semaphore, N)
// where N is the number of resources available.
type Semaphore chan empty

// P is the wait construct
// change n to consume n resources
func (s Sempahore) P(n int) {
	e := empty{}
	for i := 0; i < n; i++ {
		s <- e
	}
}

// V is the signal construct of the Semaphore
// signal n resources at a time
func (s Sempahore) V(n int) {
	for i := 0; i < n; i++ {
		<- s
	}
}

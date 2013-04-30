/* Editor Settings: expandtabs and use 4 spaces for indentation
 * ex: set softtabstop=4 tabstop=8 expandtab shiftwidth=4: *
 * -*- mode: c, c-basic-offset: 4 -*- */

/*
 * Authors: Patrick Andrew <missionsix@gmail.com>
 */

package Semaphore

type empty struct {}
type Semaphore chan empty

func (s Sempahore) P(n int) {
	e := empty{}
	for i := 0; i < n; i++ {
		s <- e
	}
}

func (s Sempahore) V(n int) {
	for i := 0; i < n; i++ {
		<- s
	}
}

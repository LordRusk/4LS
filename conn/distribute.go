/* arbitrarily structured top down controlled goroutine worker hierarchy specifically designed for 4LS */
package conn

import (
	"fmt"
	// "time"
)

const (
	Employee = iota /* nobody is below this worker and an arbitrary number of workers above */
	Manager         /* below and above this worker is an arbitrary number of workers */
	Ceo             /* nobody is above this worker and an arbitrary number of workers below - Worker.rest is recursive */
)

/* two way communication */
type roundTrip struct {
	down chan struct{}
	up   chan error
}

func newRoundTrip() roundTrip {
	return roundTrip{
		down: make(chan struct{}),
		up:   make(chan error),
	}
}

func (rt roundTrip) close() {
	close(rt.down)
	close(rt.up)
}

/* Worker starts looping Worker.work immediately until .*/
type Worker struct {
	work       func(c *Con)
	c          *Con
	isEmpty    bool /* for some Ceo's */
	closeChans map[string]roundTrip
	rest       roundTrip
	Workers    map[string]*Worker
	t          int
	id         string
	closing    bool
}

func (w *Worker) CloseUnderling(name string) {
	if !w.closing {
		w.closing = true
		rt := w.closeChans[name]
		fmt.Printf("sending rt down: id=%s\n", name)
		rt.down <- struct{}{}
		// fmt.Printf("%s: sent, waiting\n", name)
		err := <-rt.up
		fmt.Printf("%s: recieved rt up\n", name)
		if err != nil {
			fmt.Printf("failed to stop worker: %s\n", err)
		}
		rt.close()
		w.closing = false
	}
}

/* concurrently closes all underlings */
func (w *Worker) CloseUnderlings() {
	fmt.Printf("inside closeunder for %s: %d underlings\n", w.id, len(w.closeChans))
	finChan := make(chan struct{})
	for name := range w.closeChans {
		go func(name string) {
			w.CloseUnderling(name)
			finChan <- struct{}{}
		}(name)
	}
	// fmt.Println("Waiting for workers to close...")
	for i := 0; i < len(w.closeChans); i++ {
		<-finChan
		// fmt.Println("worker closed")
	}
	// fmt.Println("finished closing workers")
}

/* concurrently closes all underlings except 'exclude' */
func (w *Worker) CloseUnderlingsBut(exclude ...string) {
	fmt.Printf("inside closeunder for %s: %d underlings\n", w.id, len(w.closeChans))
	finChan := make(chan struct{})
	for name := range w.closeChans {
		for _, tname := range exclude {
			if name == tname {
				fmt.Println("dont close")
				continue
			}
		}
		go func(name string) {
			w.CloseUnderling(name)
			finChan <- struct{}{}
		}(name)
		w.CloseUnderling(name)
	}
	for i := 0; i < len(w.closeChans); i++ {
		<-finChan
	}
}

/* CloseUnderling for a underling, helpful in avoiding memory shit */
func (w *Worker) CloseUnderlingsSquared(name string) {
	w.Workers[name].CloseUnderlings()
}

/* CloseUnderlingBut for a underling, helpful in avoiding memory shit */
func (w *Worker) CloseUnderlingsButSquared(name string, exclude ...string) {
	w.Workers[name].CloseUnderlingsBut(exclude[:]...)
}

func (w *Worker) Close() {
	fmt.Printf("%+v\n", w.closeChans)
	w.CloseUnderlings()
	if w.t != Ceo {
		w.rest.up <- nil
	}

	w.rest.close() /* */
}

/* t either Ceo, Manager, Employee */
func newWorker(work func(c *Con), c *Con, isWorkEmpty bool, rest roundTrip, t int, id string) *Worker {
	w := &Worker{
		work:       work,
		c:          c,
		isEmpty:    isWorkEmpty,
		closeChans: make(map[string]roundTrip),
		rest:       rest,
		Workers:    make(map[string]*Worker),
		t:          t,
		id:         id,
	}

	// go func() {
	//time.Sleep(time.Duration(5) * time.Second) /* ?? */
	go func() {
		if w.isEmpty {
			select {
			case <-w.rest.down:
				return
			}
		}

		for {
			select {
			case <-w.rest.down:
				fmt.Printf("%s: recieved rt down! closing...\n", w.id)
				// w.Close()
				w.CloseUnderlings()
				fmt.Printf("%s: closed! Sending rt up\n", w.id)
				w.rest.up <- nil
				fmt.Printf("%s: rt up recieved! returning...\n", w.id)
				return
			default:
				w.work(w.c)
			}
		}
	}()
	//}()

	return w
}

/* returns a ceo */
func NewCeo(work func(c *Con), c *Con, isEmpty bool) *Worker {
	return newWorker(work, c, isEmpty, newRoundTrip(), Ceo, "zero")
}

/* returns a new Employee, implement's above's *Con */
func (w *Worker) New(work func(c *Con), id string) {
	if w.t == Employee {
		w.t = Manager
	}

	rt := newRoundTrip()
	w.closeChans[id] = rt
	w.Workers[id] = newWorker(work, w.c, false, rt, Employee, id)
}

func (w *Worker) Type() int {
	return w.t
}

func (w *Worker) ID() string {
	return w.id
}

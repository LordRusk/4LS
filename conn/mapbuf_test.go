package conn

import (
	"fmt"
	"testing"
	// "time"
)

/* will panic if fails */
func TestIntPost(_ *testing.T) {
	m := NewEmptyIntPost()
	c := make(chan struct{}, 10)
	for i := 1; i < 10; i++ {
		ii := i
		go func() {
			m.Map(ii, Post{No: ii})
			c <- struct{}{}
		}()
	}
	for i := 1; i < 10; i++ {
		<-c
	}

	for i := 1; i < 10; i++ {
		ii := i
		go func() {
			/* fmt.Printf("%+v\n", m.Find(ii)) /* */
			m.Find(ii)
			c <- struct{}{}
		}()
	}
	for i := 1; i < 10; i++ {
		<-c
	}

	mm := m.GetMap()
	fmt.Printf("%+v\n", mm)
}

/* will panic if fails */
func TestIntST(_ *testing.T) {
	m := NewEmptyIntST()
	c := make(chan struct{}, 10)
	for i := 1; i < 10; i++ {
		ii := i
		go func() {
			m.Map(ii, SmallThread{No: ii})
			c <- struct{}{}
		}()
	}
	for i := 1; i < 10; i++ {
		<-c
	}

	for i := 1; i < 10; i++ {
		ii := i
		go func() {
			/* fmt.Printf("%+v\n", m.Find(ii)) /* */
			m.Find(ii)
			c <- struct{}{}
		}()
	}
	for i := 1; i < 10; i++ {
		<-c
	}

	mm := m.GetMap()
	fmt.Printf("%+v\n", mm)
}

/* shouldn't fail */
func TestIntPostRange(_ *testing.T) {
	m := NewEmptyIntPost()
	c := make(chan struct{}, 10)
	for i := 1; i < 10; i++ {
		ii := i
		go func() {
			m.Map(ii, Post{No: ii})
			c <- struct{}{}
		}()
	}
	for i := 1; i < 10; i++ {
		<-c
	}

	for packPost := range m.Range() {
		fmt.Printf("k: %d, v: %d\n", packPost.K, packPost.V.No)
	}

	fmt.Println("done")
}

/* shouldn't fail */
func TestIntSTRange(_ *testing.T) {
	m := NewEmptyIntST()
	c := make(chan struct{}, 10)
	for i := 1; i < 10; i++ {
		ii := i
		go func() {
			m.Map(ii, SmallThread{No: ii})
			c <- struct{}{}
		}()
	}
	for i := 1; i < 10; i++ {
		<-c
	}

	for packPost := range m.Range() {
		fmt.Printf("k: %d, v: %d\n", packPost.K, packPost.V.No)
	}
	fmt.Println("done")
}

/*
this was used to sniff out IntPost.NewUnderlyer not swapping
the underlying map causing big memeory usage, this shouldn't be ran
*/
/*
func TestIntPostMem(_ *testing.T) {
	m := NewEmptyIntPost()
	c := make(chan struct{}, 1000)
	ec := make(chan struct{})
	go func() {
		i := 0
		for time.Now().Before(time.Now().Add(time.Duration(1) * time.Minute)) {
			i++
			println(m.Len())
			if m.Len() > 10000 {
				m.NewUnderlyer(make(map[int]Post))
				m.Map(i, Post{No: i})
				fmt.Printf("post newunderlyer find: %+v\n", m.Find(i))
				fmt.Println(m.Len())
			} else if m.Len() == 0 {
				// println("why 0")
				// if i > 10 {
					// return
				// }
			}
			m.Map(i, Post{No: i})
			c <- struct{}{}
		}
		ec <- struct{}{}
	}()
	for {
		select {
		case <-c:
		case <-ec:
			break
		}
	}
}
*/

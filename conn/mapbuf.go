/*
bufmap allows for concurrent read/writes to
map[int]conn.Post and map[int]conn.SmallThread
*/
package conn

import (
	"fmt"
	"time"
)

const buflength = 4012 * 4 /* arbitrarily large */

const (
	/* instructions */
	Map = iota
	Find
	Remove
	Range
	Change
)

type PackPost struct {
	instruction int
	K           int
	V           Post
}

type PackST struct {
	instruction int
	K           int
	V           SmallThread
}

type postWrap struct {
	post        Post
	rangeKeys   []int
	rangeValues []Post
}

/* ported from 4LS/distribute.roundTrip */
type roundTripPost struct {
	down chan PackPost /* key and value */
	up   chan postWrap
}

func newRoundTripPost() roundTripPost {
	return roundTripPost{
		down: make(chan PackPost, buflength),
		up:   make(chan postWrap, buflength),
	}
}

type smallThreadWrap struct {
	thread      SmallThread
	rangeKeys   []int
	rangeValues []SmallThread
}

type roundTripST struct {
	down chan PackST /* key and value */
	up   chan smallThreadWrap
}

func newRoundTripST() roundTripST {
	return roundTripST{
		down: make(chan PackST, buflength),
		up:   make(chan smallThreadWrap, buflength),
	}
}

type emptyTrip struct {
	down, up chan struct{}
}

func newEmptyTrip() emptyTrip {
	return emptyTrip{
		down: make(chan struct{}),
		up:   make(chan struct{}),
	}
}

/* conncurent map[int]conn.Post */
type IntPost struct {
	m    map[int]Post
	buf  roundTripPost
	kill emptyTrip
}

func NewIntPost(mp map[int]Post) *IntPost {
	m := IntPost{
		m:    mp,
		buf:  newRoundTripPost(),
		kill: newEmptyTrip(),
	}

	/* map, find, remove */
	go func() {
		handlePack := func(p PackPost) {
			switch p.instruction {
			case Map:
				m.m[p.K] = p.V
				m.buf.up <- postWrap{}
			case Find:
				m.buf.up <- postWrap{post: m.m[p.K]}
			case Remove:
				delete(m.m, p.K)
				m.buf.up <- postWrap{}
			case Range:
				pw := postWrap{}
				for k, v := range m.m {
					pw.rangeKeys = append(pw.rangeKeys, k)
					pw.rangeValues = append(pw.rangeValues, v)
				}
				m.buf.up <- pw
			default:
				fmt.Println("bufmap just default'd on instruction....how is pack.key > 4 || < 0?")
			}
		}

		/* avoid write on nil map */
		time.Sleep(2)

		for {
			select {
			case <-m.kill.down:
				for i := 0; i < len(m.buf.down); i++ {
					pack := <-m.buf.down
					handlePack(pack)
				}
				m.kill.up <- struct{}{}
				return
			case pack := <-m.buf.down:
				handlePack(pack)
			}
		}
	}()

	return &m
}

func NewEmptyIntPost() *IntPost {
	return NewIntPost(make(map[int]Post))
}

/*
replacing underlying map with mm
make sure no conncurent actions are being made
*/
func (m *IntPost) NewUnderlyer(mm map[int]Post) {
	m.m = mm
}

/* returns underlying map */
func (m *IntPost) GetMap() map[int]Post {
	return m.m
}

func (m *IntPost) Clean() {
	m.NewUnderlyer(make(map[int]Post))
}

/* closes the map buffer, all buffered map updates are handled before Close() finishes */
func (m *IntPost) Close() map[int]Post {
	m.kill.down <- struct{}{}
	<-m.kill.up
	mm := m.GetMap()
	m.Clean()
	return mm
}

/* add a new entry */
func (m *IntPost) Map(key int, value Post) {
	m.buf.down <- PackPost{K: key, V: value, instruction: Map}
	<-m.buf.up
}

/* get a value from a key */
func (m *IntPost) Find(key int) Post {
	m.buf.down <- PackPost{K: key, V: Post{}, instruction: Find}
	pw := <-m.buf.up
	return pw.post
}

/* same as delete(map, key) */
func (m *IntPost) Remove(key int) {
	m.buf.down <- PackPost{K: key, V: Post{}, instruction: Remove}
	<-m.buf.up
}

/* returns a dirty answer */
func (m *IntPost) Len() int {
	return len(m.m)
}

/*
returns a channel that will send all keys and values,
closing after, effectively replacing `k, v := range map`
*/
func (m *IntPost) Range() chan PackPost {
	m.buf.down <- PackPost{K: 0, V: Post{}, instruction: Range}
	pw := <-m.buf.up
	ch := make(chan PackPost, len(pw.rangeKeys))
	for i := 0; i < len(pw.rangeKeys); i++ {
		ch <- PackPost{K: pw.rangeKeys[i], V: pw.rangeValues[i]}
	}
	close(ch)
	return ch
}

/* conncurent map[int]conn.SmallThread */
type IntST struct {
	m    map[int]SmallThread
	buf  roundTripST
	kill emptyTrip
}

func NewIntST(mp map[int]SmallThread) *IntST {
	m := IntST{
		m:    mp,
		buf:  newRoundTripST(),
		kill: newEmptyTrip(),
	}

	/* map, find, remove */
	go func() {
		handlePack := func(p PackST) {
			switch p.instruction {
			case Map:
				m.m[p.K] = p.V
				m.buf.up <- smallThreadWrap{}
			case Find:
				m.buf.up <- smallThreadWrap{thread: m.m[p.K]}
			case Range:
				stw := smallThreadWrap{}
				for k, v := range m.m {
					stw.rangeKeys = append(stw.rangeKeys, k)
					stw.rangeValues = append(stw.rangeValues, v)
				}
				m.buf.up <- stw
			default:
				fmt.Println("bufmap just default'd on instruction....how is pack.key > 4 || < 0?")
			}
		}

		for {
			select {
			case <-m.kill.down:
				for i := 0; i < len(m.buf.down); i++ {
					pack := <-m.buf.down
					handlePack(pack)
				}
				m.kill.up <- struct{}{}
				return
			case pack := <-m.buf.down:
				handlePack(pack)
			}
		}
	}()

	return &m
}

func NewEmptyIntST() *IntST {
	return NewIntST(make(map[int]SmallThread))
}

/*
replacing underlying map with mm
make sure no conncurent actions are being made
*/
func (m *IntST) NewUnderlyer(mm map[int]SmallThread) {
	//fmt.Printf("Just set underlying map to %+v\n", mm)
	m.m = mm
	// fmt.Printf("underlying map: %+v\n", m.m)
}

/* returns underlying map */
func (m *IntST) GetMap() map[int]SmallThread {
	return m.m
}

func (m *IntST) Clean() {
	m.NewUnderlyer(make(map[int]SmallThread)) /* garbage collector should take care of the rest */
}

/* closes the map buffer, all buffered map updates are handled before Close() finishes */
func (m *IntST) Close() map[int]SmallThread {
	// fmt.Println("Finishings buffered instructions")
	m.kill.down <- struct{}{}
	<-m.kill.up
	// fmt.Println("finished buffered instructions")
	mm := m.GetMap()
	// fmt.Println("Got Internal Map, cleaning...")
	m.Clean()
	return mm
}

/* add a new entry */
func (m *IntST) Map(key int, value SmallThread) {
	m.buf.down <- PackST{K: key, V: value, instruction: Map}
	<-m.buf.up
}

/* get a value from a key */
func (m *IntST) Find(key int) SmallThread {
	m.buf.down <- PackST{K: key, V: SmallThread{}, instruction: Find}
	stw := <-m.buf.up
	return stw.thread
}

/* same as delete(map, key) */
func (m *IntST) Remove(key int) {
	m.buf.down <- PackST{K: key, V: SmallThread{}, instruction: Remove}
	<-m.buf.up
}

/* returns a dirty answer and stops working after new underlyer */
func (m *IntST) Len() int {
	return len(m.m)
}

/*
returns a channel that will send all keys and values,
closing after, effectively replacing `k, v :=range map`
*/
func (m *IntST) Range() chan PackST {
	m.buf.down <- PackST{K: 0, V: SmallThread{}, instruction: Range}
	stw := <-m.buf.up
	ch := make(chan PackST, len(stw.rangeKeys))
	for i := 0; i < len(stw.rangeKeys); i++ {
		ch <- PackST{K: stw.rangeKeys[i], V: stw.rangeValues[i]}
	}
	close(ch)
	return ch
}

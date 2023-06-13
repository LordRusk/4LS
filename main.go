package main

import (
	// "encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	// errs "github.com/pkg/errors"
	"github.com/pkg/profile"

	// "github.com/lordrusk/4LS/bord"
	"github.com/lordrusk/4LS/conn"
)

var logFolder = flag.String("folder", "~/fourall", "Set the folder to read/write internal structues to")
var skipF = flag.Bool("skip", false, "Skip reading from -folder")
var waitTime = flag.Int("wait", 15, "How long to wait before sweeping")
var saveTime = flag.Int("time", 60, "How many minutes between saving to -folder")

// profiling
var cpuProf = flag.Bool("cpu", false, "Create a cpu profile (Overwrites -mem)")
var memProf = flag.Bool("mem", false, "Create a memory profile")
var memProfRate = flag.Int("memRate", 0, "Set the memory profile rate")

var imgBuf = 8112

/* wait's for CTRL+C */
func waitCtrlC(s *bool) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	_ = <-c

	ss := *s
	ss = true
	s = &ss
}

func main() {
	flag.Parse()

	/* profiling */
	if *cpuProf {
		defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	} else if *memProf {
		if *memProfRate != 0 {
			defer profile.Start(profile.MemProfile, profile.MemProfileRate(*memProfRate), profile.ProfilePath(".")).Stop()
		} else {
			defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()
		}
	}

	c := conn.NewCon(*logFolder)
	defer func() {
		if err := c.Close(); err != nil { /* c.(*logger.Logger).Close() */
			fmt.Printf("Failed to close client: %s\n", err)
		}
	}()
	bts := []conn.Bored{ /* boards to scrape */
		c.X,
		c.POL,
		c.K,
		c.B,
	}

	if !*skipF {
		for _, v := range c.Map() {
			if err := c.GetBoard(v); err != nil {
				c.Printf("%s\n", err)
			}
		}
	}

	/* image downloads */
	imgChan := make(chan conn.Image, imgBuf)
	path := strings.ReplaceAll(*logFolder, "~", "$HOME")
	path = os.ExpandEnv(path)
	for _, v := range c.Map() {
		if err := os.MkdirAll(path+"/"+v.Name, os.ModePerm); err != nil {
			fmt.Printf("%s: Cannot create '%s'", err, path)
		}
	}

	// fumc := func() {
	// }

	/* main work loop */
	lastUpdate, _ := time.Parse(time.Layout, time.Layout)
	firstRun := true
	ctrlC := false
	funk := func(c *conn.Con) {
		if lastUpdate.Add(time.Duration(*saveTime) * time.Minute).Before(time.Now()) {
			lastUpdate = time.Now()
			/* this is ran once at the beginning */
			if firstRun {
				firstRun = false
			} else {
				fmt.Println("Saving and cleaning...")
				/* underlying map isn't being cleaned and gigabytes of memory are being used.
				 * Fix This before implementing image downloads and a CLI interface */

				/* this shit is still happening, signs say its within mapbuf, but I can't read the memory profile or something. */
				// c.w.workers["board master"].CloseUnderlingsBut("imgcollecter")
				c.W.CloseUnderlingsSquared("Board Master")
				/* fmt.Println("just closed but ") /* */
				fmt.Println("just closed ")
				for _, v := range c.Map() { /* */
					if err := c.PutBoard(v); err != nil {
						c.Printf("%s\n", err)
					}
				}
				fmt.Println("just put boards ")

				for _, v := range c.Map() { /* */
					if err := c.GetBoard(v); err != nil {
						c.Printf("%s\n", err)
					}
				}

				fmt.Println("just got boards ")

			}

			for _, board := range bts {
				b := board
				funky := func(c *conn.Con) {
					//fmt.Printf("Sweeping... %s\n", b.Name)
					if err := sweep(c, &b, imgChan); err != nil {
						c.Printf("%s\n", err)
						return
					}
					//fmt.Printf("Swept %s\n", b.Name)

					for i := *waitTime; i > -1; i-- {
						// fmt.Printf("sleeping... %s\n", b.Name)
						time.Sleep(1 * time.Second)
					}

				}
				c.W.Workers["Board Master"].New(funky, board.Name)
			}

		}
		time.Sleep(time.Duration(*saveTime/3) * time.Minute)

	}
	c.W.New(funk, "Board Master")

	waitCtrlC(&ctrlC)
	c.W.CloseUnderlings()

	for _, v := range c.Map() { /* */
		if err := c.PutBoard(v); err != nil {
			c.Printf("%s\n", err)
		}
	}
}

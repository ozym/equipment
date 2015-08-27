package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/ozym/dmc"
	"github.com/ozym/zone"
)

func status(args []string) {

	f := flag.NewFlagSet("status", flag.ExitOnError)
	f.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Check for any changes in installed equipment state\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options] status [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "General Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Equipment Status Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		f.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	var master string
	flag.StringVar(&master, "master", "rhubarb.geonet.org.nz.", "default master for equipment service lookup")

	var lookup string
	flag.StringVar(&lookup, "zone", "wan.geonet.org.nz.", "default zone for equipment service lookup")

	var timeout time.Duration
	f.DurationVar(&timeout, "timeout", time.Second*10, "provide a service timeout")

	var retries int
	f.IntVar(&retries, "retries", 3, "provide a service retry")

	var limit int
	f.IntVar(&limit, "limit", 8, "number of concurrent queries")

	var models string
	f.StringVar(&models, "models", ".*", "regex expression to match equipment models")

	var sites string
	f.StringVar(&sites, "sites", ".*", "regex expression to match equipment sites")

	if err := f.Parse(args); err != nil {
		f.Usage()

		log.Fatalf("Invalid option(s) given")
	}

	m := regexp.MustCompile(models)
	s := regexp.MustCompile(sites)

	// concurrent goroutines
	var wg sync.WaitGroup

	// semaphore to limit number of goroutines
	sem := make(chan struct{}, limit)

	details, err := zone.LoadLocal(master, []string{lookup}, []string{})
	if err != nil {
		log.Fatal(err)
	}

	devices := details.MustMatchByModel(models).MustMatchByName(sites)

	for _, l := range devices.List {
		if !m.Match(([]byte)(l.Model)) {
			continue
		}
		if !s.Match(([]byte)(l.Name)) {
			continue
		}

		d := dmc.Device{
			Name:  l.Name,
			IP:    l.IP,
			Model: l.Model,
		}

		if verbose {
			log.Printf("checking equipment: %s\n", d.String())
		}

		for _, m := range dmc.ModelList {
			if !d.Match(m) {
				continue
			}

			// wait ...
			sem <- struct{}{}
			wg.Add(1)

			go func(m dmc.Model, d dmc.Device) {
				defer func() { <-sem; wg.Done() }()

				if verbose {
					log.Printf("checking: %s against %s\n", d.String(), m.Name())
				}

				if s, _ := d.Identify(m, d.Model, timeout, retries); s != nil {
					if device(d, s) {
						return
					}
				}
				if s := d.Discover(m, d.Model, timeout, retries); s != nil {
					if device(d, s) {
						return
					}
				}

				if verbose {
					log.Printf("missed equipment: %s\n", d.String())
				}
			}(m, d)
		}
	}

	wg.Wait()
}

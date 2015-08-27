package main

import (
	"sync"

	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ozym/dmc"
	"github.com/ozym/zone"
)

func check(args []string) {

	f := flag.NewFlagSet("check", flag.ExitOnError)
	f.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Check for any newly installed equipment (tagged as uninstalled)\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options] check [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "General Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "New Equipment Options:\n")
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
	f.StringVar(&models, "models", "^Uninstalled", "device model regexp to check")

	var sites string
	f.StringVar(&sites, "sites", ".*", "device name regexp to check")

	if err := f.Parse(args); err != nil {
		f.Usage()

		log.Fatalf("Invalid option(s) given")
	}

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
		model := strings.Replace(l.Model, "Uninstalled ", "", -1)

		d := dmc.Device{
			Name:  l.Name,
			IP:    l.IP,
			Model: model,
		}

		if verbose {
			log.Printf("checking: %s\n", d.String())
		}

		for _, m := range dmc.ModelList {

			// wait ...
			sem <- struct{}{}
			wg.Add(1)

			go func(m dmc.Model, d dmc.Device) {
				defer func() { <-sem; wg.Done() }()

				if verbose {
					log.Println(d)
				}

				if s, _ := d.Identify(m, model, timeout, retries); s != nil {
					if device(d, s) {
						return
					}
				}

				if s := d.Discover(m, model, timeout, retries); s != nil {
					if device(d, s) {
						return
					}
				}
			}(m, d)
		}
	}

	wg.Wait()
}

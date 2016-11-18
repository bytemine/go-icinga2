package icinga2

import (
	"encoding/json"
	"github.com/bytemine/go-icinga2/event"
	"log"
	"sync"
	"testing"
	"time"
)

// ExampleVagrant connects to a icinga2 vagrant box, found at
// https://github.com/Icinga/icinga-vagrant.git
// Use icinga2x and add port 5665 to the forwarded ports.
func ExampleVagrant(t *testing.T) {
	i, err := NewClient("https://localhost:5665", "root", "icinga", true)
	if err != nil {
		t.Fatal(err)
	}

	// Get an event stream for StateChange and CheckResult events.
	es, err := i.EventStream("testing", "", event.StreamTypeStateChange, event.StreamTypeCheckResult)
	if err != nil {
		t.Fatal(err)
	}

	// Multiplex StateChange and CheckResult events into seperate readers.
	m := event.NewMux(es, event.StreamTypeStateChange, event.StreamTypeCheckResult)

	// Get the reader for StateChange events.
	st, err := m.Reader(event.StreamTypeStateChange)
	if err != nil {
		t.Fatal(err)
	}

	// Get the reader for CheckResult events.
	cr, err := m.Reader(event.StreamTypeCheckResult)
	if err != nil {
		t.Fatal(err)
	}

	// Receive the events with seperate goroutines.
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		decx := json.NewDecoder(st)

		var x event.StateChange

		for {
			err := decx.Decode(&x)
			if err != nil {
				return
			}
			log.Printf("%#v", x)
		}

		return
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		decy := json.NewDecoder(cr)

		var y event.CheckResult
		for {
			err := decy.Decode(&y)
			if err != nil {
				return
			}
			log.Printf("%#v", y)
		}
		return
	}()

	go func() {
		time.Sleep(time.Second * 5)
		m.Close()
	}()

	wg.Wait()
}

package event

import (
	"encoding/json"
	"io"
	"sync"
	"testing"
)

var testEvents = []string{
	`{"type":"StateChange","timestamp":1234.5678,"host":"test.example.org","service":"testservice","state":1.0,"state_type":1.0}`,
	`{"type":"StateChange","timestamp":1235.5678,"host":"test.example.org","service":"testservice","state":0.0,"state_type":1.0}`,
	`{"type":"CheckResult","timestamp":1236.5678,"host":"test.example.org","service":"testservice"}`,
	`{"type":"StateChange","timestamp":1237.5678,"host":"test.example.org","service":"testservice","state":0.0,"state_type":1.0}`,
	`{"type":"CheckResult","timestamp":1238.5678,"host":"test.example.org","service":"testservice"}`,
	`{"type":"StateChange","timestamp":1235.5678,"host":"test.example.org","service":"testservice","state":0.0,"state_type":1.0}`,
}

func testEventReader(t *testing.T, events []string) (io.Reader) {
	rp, wp := io.Pipe()
	go func() {
		for _, v := range testEvents {
			_, err := wp.Write(append([]byte(v), '\n'))
			if err != nil {
				t.Fatal(err)
			}
		}
		wp.Close()
	}()

	return rp
}

func TestMux(t *testing.T) {
	buf := testEventReader(t, testEvents)

	m := NewMux(buf, StreamTypeStateChange, StreamTypeCheckResult)
	st, err := m.Reader(StreamTypeStateChange)
	if err != nil {
		t.Fatal(err)
	}
	cr, err := m.Reader(StreamTypeCheckResult)
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		dec := json.NewDecoder(st)

		var x StateChange

		for dec.More() {
			err := dec.Decode(&x)
			if err != nil {
				t.Log(err)
				return
			}
			t.Logf("%#v", x)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		dec := json.NewDecoder(cr)

		var y CheckResult

		for dec.More() {
			err := dec.Decode(&y)
			if err != nil {
				t.Log(err)
				return
			}
			t.Logf("%#v", y)
		}
	}()

	wg.Wait()
}

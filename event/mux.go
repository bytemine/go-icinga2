package event

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"sync"
)

type StreamType string

// Event stream selectors, see: http://docs.icinga.org/icinga2/snapshot/doc/module/icinga2/chapter/icinga2-api#icinga2-api-event-streams
const (
	StreamTypeCheckResult            StreamType = "CheckResult"
	StreamTypeStateChange            StreamType = "StateChange"
	StreamTypeNotification           StreamType = "Notification"
	StreamTypeAcknowledgementSet     StreamType = "AcknowledgementSet"
	StreamTypeAcknowledgementCleared StreamType = "AcknowledgementCleared"
	StreamTypeCommentAdded           StreamType = "CommentAdded"
	StreamTypeCommentRemoved         StreamType = "CommentRemoved"
	StreamTypeDowntimeAdded          StreamType = "DowntimeAdded"
	StreamTypeDowntimeRemoved        StreamType = "DowntimeRemoved"
	StreamTypeDowntimeTriggered      StreamType = "DowntimeTriggered"
)

// Mux multiplexes a single icinga event stream reader into multiple readers for different StreamTypes.
type Mux struct {
	mux  map[StreamType]pipe
	muxm sync.RWMutex
}

type pipe struct {
	r *io.PipeReader
	w *io.PipeWriter
}

// NewMux initializes a new Mux. The StreamTypes used here must be consumed by aquiring io.Readers
// using the Reader method and reading from them, otherwise reception of events will stall.
func NewMux(r io.Reader, types ...StreamType) *Mux {
	var m Mux
	m.muxm.Lock()
	m.mux = make(map[StreamType]pipe)
	for _, v := range types {
		r, w := io.Pipe()
		m.mux[v] = pipe{r, w}
	}
	m.muxm.Unlock()

	go func() {
		br := bufio.NewReader(r)

		for {
			buf, err := br.ReadBytes('\n')
			if err != nil {
				m.closeWithError(err)
				return
			}

			x := struct {
				Type StreamType
			}{}

			err = json.Unmarshal(buf, &x)
			if err != nil {
				m.closeWithError(fmt.Errorf("Unregistered stream type %v", x.Type))
				return
			}

			out, ok := m.mux[x.Type]
			if !ok {
				if err != nil {
					m.closeWithError(fmt.Errorf("Unregistered stream type %v", x.Type))
					return
				}
			}

			_, err = out.w.Write(append(buf, '\n'))
			if err != nil {
				m.closeWithError(err)
				return
			}
		}
	}()

	return &m
}

// Reader returns a io.Reader for the specified StreamType, usable with a
// json.Decoder. If the StreamType wasn't registered for multiplexing with
// the Mux an error is returned.
func (m *Mux) Reader(typ StreamType) (io.Reader, error) {
	m.muxm.RLock()
	p, ok := m.mux[typ]
	m.muxm.RUnlock()

	if !ok {
		return nil, fmt.Errorf("Unregistered stream type %v", typ)
	}

	return p.r, nil
}

func (m *Mux) closeWithError(err error) error {
	m.muxm.RLock()
	for _, v := range m.mux {
		err := v.w.CloseWithError(err)
		if err != nil {
			return err
		}
	}
	m.muxm.RUnlock()
	return nil
}

// Close closes all Readers.
func (m *Mux) Close() error {
	m.muxm.RLock()
	for _, v := range m.mux {
		err := v.w.Close()
		if err != nil {
			return err
		}
	}
	m.muxm.RUnlock()
	return nil
}

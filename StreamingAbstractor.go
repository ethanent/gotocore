package protocore

import (
	"errors"
	"sync"
)

// HandlerChan is a channel capable of receiving data parsed using a schema
type HandlerChan chan map[string]interface{}

// StreamingAbstractor is a streaming agent for transmitting data
type StreamingAbstractor struct {
	schemas     map[string]Schema
	outMux      *sync.RWMutex
	outBuffer   []byte
	inMux       *sync.Mutex
	inBuffer    []byte
	handlers    map[string][]HandlerChan
	frameSchema *Schema
}

// NewStreamingAbstractor initializes and returns a StreamingAbstractor pointer
func NewStreamingAbstractor() *StreamingAbstractor {
	frameSch := Schema{}

	frameSch.Components = append(frameSch.Components, Component{
		Name: "event",
		Kind: String,
	})

	frameSch.Components = append(frameSch.Components, Component{
		Name: "mode",
		Kind: UInt,
		Size: 8,
	})

	frameSch.Components = append(frameSch.Components, Component{
		Name: "serialized",
		Kind: Buffer,
		Size: 8,
	})

	return &StreamingAbstractor{
		schemas:     map[string]Schema{},
		outMux:      &sync.RWMutex{},
		outBuffer:   []byte{},
		inBuffer:    []byte{},
		frameSchema: &frameSch,
	}
}

// Register saves a schema into StreamingAbstractor s
func (s *StreamingAbstractor) Register(name string, schema Schema) {
	s.schemas[name] = schema
}

// Send sends data from StreamingAbstractor s
func (s *StreamingAbstractor) Send(name string, data map[string]interface{}) error {
	relevSchema, ok := s.schemas[name]

	if !ok {
		return errors.New("unregistered schema '" + name + "'")
	}

	d, err := relevSchema.Build(data)

	if err != nil {
		return err
	}

	s.outMux.Lock()
	s.outBuffer = append(s.outBuffer, d...)
	s.outMux.Unlock()

	return nil
}

// Handle instructs StreamingAbstractor s to push to the channel once a message of name name is received
func (s *StreamingAbstractor) Handle(name string, ch HandlerChan) {
	_, ok := s.schemas[name]

	if !ok {
		panic("Unregistered schema '" + name + "'")
	}

	schHandlers, ok := s.handlers[name]

	if !ok {
		s.handlers[name] = []HandlerChan{ch}
	} else {
		schHandlers = append(schHandlers, ch)
	}
}

// Handle io functionality

func (s *StreamingAbstractor) Read(p []byte) (int, error) {
	s.outMux.RLock()
	defer s.outMux.RUnlock()

	readCount := 0

	for readCount < len(s.outBuffer)-1 && readCount < len(p)-1 {
		p = append(p, s.outBuffer[readCount])

		readCount++
	}

	s.outBuffer = s.outBuffer[readCount:]

	return readCount, nil
}

func (s *StreamingAbstractor) Write(p []byte) (int, error) {
	s.inMux.Lock()
	defer s.inMux.Unlock()

	s.inBuffer = append(s.inBuffer, p...)

	// Attempt parse frame

	return len(p), nil
}

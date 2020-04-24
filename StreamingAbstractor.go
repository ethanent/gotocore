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
		handlers:    map[string][]HandlerChan{},
		schemas:     map[string]Schema{},
		outMux:      &sync.RWMutex{},
		outBuffer:   []byte{},
		inMux:       &sync.Mutex{},
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

	serialized, err := relevSchema.Build(data)

	if err != nil {
		return err
	}

	fullFrame, err := s.frameSchema.Build(map[string]interface{}{
		"event":      name,
		"mode":       uint(0),
		"serialized": serialized,
	})

	if err != nil {
		return err
	}

	s.outMux.Lock()
	defer s.outMux.Unlock()
	s.outBuffer = append(s.outBuffer, fullFrame...)

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
		s.handlers[name] = schHandlers
	}
}

// Handle io functionality

func (s *StreamingAbstractor) Read(p []byte) (int, error) {
	s.outMux.RLock()
	defer s.outMux.RUnlock()

	readCount := 0

	for readCount < len(s.outBuffer) && readCount < len(p) {
		p[readCount] = s.outBuffer[readCount]

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

	frameData, frameLen, frameErr := s.frameSchema.Parse(p)

	if frameErr == nil {
		// Frame is now fully arrived

		// Clear frame from buffer

		s.inBuffer = s.inBuffer[frameLen:]

		// Handle frame data

		eventName := frameData["event"].(string)

		relevSchema, ok := s.schemas[eventName]

		if ok {
			// Schema is relevant

			data, _, err := relevSchema.Parse(frameData["serialized"].([]byte))

			if err == nil {
				// Find handlers and give them data, if any exist

				s.informHandlers(eventName, data)
			} else {
				// Malformatted serialized data. Will be ignored.
			}
		} else {
			// Schema is unregistered. Will ignore message.
		}
	} else {
		// Frame is not yet fully buffered.
	}

	return len(p), nil
}

func (s *StreamingAbstractor) informHandlers(event string, data map[string]interface{}) {
	evHandlers, ok := s.handlers[event]

	if !ok {
		// No handlers for event

		return
	}

	// Handlers exist for event!
	for len(evHandlers) > 0 {
		evh := evHandlers[0]

		// Provide data to handler
		evh <- data

		// Drop handler
		evHandlers = evHandlers[1:]
	}

	s.handlers[event] = evHandlers
}

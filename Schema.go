package protocore

// ComponentType is an enum representing a Component's type
type ComponentType int

const (
	// Varint is a variable length integer
	Varint ComponentType = iota

	// Buffer is an []byte
	Buffer
)

// Component is an element of a protocol
type Component struct {
	Name string
	Kind ComponentType
}

// Schema is a Protocore schema to be used for encoding or decoding
type Schema struct {
	Components []Component
}

// Parse parses []byte buf into a map[string]interface{}
func (s *Schema) Parse(buf []byte) (data map[string]interface{}, endRead int, err error) {
	curIdx := 0
	build := map[string]interface{}{}

	for _, curComponent := range s.Components {
		switch curComponent.Kind {
		case Varint:
			val, readBytes, err := parseVarint(buf, curIdx, &curComponent)

			if err != nil {
				return nil, 0, err
			}

			build[curComponent.Name] = val
			curIdx += readBytes
		case Buffer:
			val, readBytes, err := parseBuffer(buf, curIdx, &curComponent)

			if err != nil {
				return nil, 0, err
			}

			build[curComponent.Name] = val
			curIdx += readBytes
		}
	}

	return build, curIdx, nil
}

// Build builds a map[string]interface{} into an []byte
func (s *Schema) Build(data map[string]interface{}) []byte {
	build := []byte{}

	for _, curComponent := range s.Components {
		switch curComponent.Kind {
		case Varint:
			build = append(build, buildVarint(data[curComponent.Name].(int))...)
		case Buffer:
			build = append(build, buildBuffer(data[curComponent.Name].([]byte))...)
		default:
			panic("Unexpected component kind.")
		}
	}

	return build
}

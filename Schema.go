package protocore

// ComponentType is an enum representing a Component's type
type ComponentType int

const (
	// Varint is a variable length integer
	Varint ComponentType = iota
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
func (s *Schema) Parse(buf []byte) (data map[string]interface{}, endByte int, err error) {
	curIdx := 0
	build := map[string]interface{}{}

	for _, curComponent := range s.Components {
		if curComponent.Kind == Varint {
			val, readBytes, err := parseVarint(buf, curIdx, &curComponent)

			if err != nil {
				return nil, 0, err
			}

			build[curComponent.Name] = val
			curIdx += readBytes
		}
	}

	return build, curIdx, nil
}

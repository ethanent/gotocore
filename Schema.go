package gotocore

import "errors"

// ComponentKind is an enum representing a Component's Kind
type ComponentKind int

const (
	// Varint is a variable length integer
	Varint ComponentKind = iota

	// Buffer is an []byte
	Buffer

	// String is a string
	String

	// Uint is an unsigned integer
	Uint
)

// Component is an element of a protocol
type Component struct {
	Name string
	Kind ComponentKind

	// Size is only used for UInt and Int kinds
	Size int
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
		case String:
			val, readBytes, err := parseString(buf, curIdx, &curComponent)

			if err != nil {
				return nil, 0, err
			}

			build[curComponent.Name] = val
			curIdx += readBytes
		case Uint:
			val, readBytes, err := parseUInt(buf, curIdx, &curComponent)

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
func (s *Schema) Build(data map[string]interface{}) ([]byte, error) {
	build := []byte{}

	for _, curComponent := range s.Components {
		switch curComponent.Kind {
		case Varint:
			dassert, ok := data[curComponent.Name].(int)

			if !ok {
				return nil, errors.New("failed to assert to int")
			}

			build = append(build, buildVarint(dassert)...)
		case Buffer:
			dassert, ok := data[curComponent.Name].([]byte)

			if !ok {
				return nil, errors.New("failed to assert to []byte")
			}

			build = append(build, buildBuffer(dassert)...)
		case String:
			dassert, ok := data[curComponent.Name].(string)

			if !ok {
				return nil, errors.New("failed to assert to string")
			}

			build = append(build, buildString(dassert)...)
		case Uint:
			dassert, ok := data[curComponent.Name].(uint)

			if !ok {
				return nil, errors.New("failed to assert to uint")
			}

			builtUint, err := buildUInt(dassert, curComponent.Size)

			if err != nil {
				return nil, err
			}

			build = append(build, builtUint...)
		default:
			panic("Unexpected component kind.")
		}
	}

	return build, nil
}

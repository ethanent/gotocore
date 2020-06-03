package gotocore

import "errors"

func parseBuffer(buf []byte, startIdx int, curComponent *Component) (value []byte, readBytes int, err error) {
	bufLen, preLen, err := parseVarint(buf, startIdx, curComponent)

	if err != nil {
		return nil, 0, err
	}

	if len(buf)-startIdx < preLen+bufLen {
		return nil, 0, errors.New("Could not complete parse. Incomplete buffer '" + curComponent.Name + "'.")
	}

	return buf[startIdx+preLen : startIdx+preLen+bufLen], preLen + bufLen, nil
}

func buildBuffer(value []byte) []byte {
	v := buildVarint(len(value))

	v = append(v, value...)

	return v
}

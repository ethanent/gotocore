package gotocore

import "errors"

func parseString(buf []byte, startIdx int, curComponent *Component) (value string, readBytes int, err error) {
	sbuf, slen, err := parseBuffer(buf, startIdx, curComponent)

	if err != nil {
		return "", 0, err
	}

	return string(sbuf), slen, nil
}

func buildString(ov interface{}) ([]byte, error) {
	value, ok := ov.(string)

	if !ok {
		return nil, errors.New("expected string value for ov")
	}

	sbuf := []byte(value)

	v := buildVarint(len(sbuf))

	v = append(v, sbuf...)

	return v, nil
}

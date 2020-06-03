package gotocore

func parseString(buf []byte, startIdx int, curComponent *Component) (value string, readBytes int, err error) {
	sbuf, slen, err := parseBuffer(buf, startIdx, curComponent)

	if err != nil {
		return "", 0, err
	}

	return string(sbuf), slen, nil
}

func buildString(value string) []byte {
	sbuf := []byte(value)

	v := buildVarint(len(sbuf))

	v = append(v, sbuf...)

	return v
}

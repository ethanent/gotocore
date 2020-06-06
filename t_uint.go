package gotocore

import (
	"errors"
	"math"
	"strconv"
)

func parseUInt(buf []byte, startIdx int, curComponent *Component) (value uint, readBytes int, err error) {
	uintLen := curComponent.Size / 8

	if len(buf)-startIdx < uintLen {
		return 0, 0, errors.New("Could not complete parse. Incomplete UInt '" + curComponent.Name + "'.")
	}

	pv := 1
	var val uint = 0

	for i := startIdx; i < startIdx+uintLen; i++ {
		val += uint(buf[i]) * uint(pv)

		pv *= 256
	}

	return val, uintLen, nil
}

func buildUInt(value uint, size int) []byte {
	if size < 1 {
		panic("invalid uint size " + strconv.Itoa(size))
	}

	if value == 0 {
		return []byte{0}
	}

	build := []byte{}

	var maxPlace uint = 1

	for value > maxPlace {
		maxPlace *= 256
	}

	maxPlace /= 256

	for maxPlace >= 1 {
		placeValue := uint(math.Floor(float64(value) / float64(maxPlace)))
		build = append([]byte{byte(placeValue)}, build...)

		value %= maxPlace
		maxPlace /= 256
	}

	// Ensure size confirmity

	if len(build)*8 > size {
		panic("uint value " + strconv.Itoa(int(value)) + " out of range for uint" + strconv.Itoa(size))
	}

	for len(build)*8 < size {
		build = append([]byte{0}, build...)
	}

	return build
}

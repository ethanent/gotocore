package protocore

import (
	"errors"
	"math"
)

func parseVarint(buf []byte, startIdx int, curComponent *Component) (value int, readBytes int, err error) {
	// Interpret prefix

	viPrefix := int(buf[startIdx])
	viLen := 0
	sign := 1 // 1 = positive, -1 = negative

	if viPrefix > 127 {
		viLen = viPrefix - 127
		sign = -1
	} else {
		viLen = viPrefix + 1
		sign = 1
	}

	if len(buf) < startIdx+viLen+1 {
		return 0, 0, errors.New("Could not complete parse. Incomplete varint '" + curComponent.Name + "'.")
	}

	viRaw := buf[startIdx+1 : startIdx+viLen+1]
	viVal := 0

	for viIdx, viCurVal := range viRaw {
		placeVal := int(math.Pow(256, float64(viIdx)))

		viVal += placeVal * int(viCurVal)
	}

	viVal *= sign

	// Return read data

	return viVal, viLen + 1, nil
}

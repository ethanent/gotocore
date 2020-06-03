package gotocore

import (
	"errors"
	"math"
)

type varintPrefixData struct {
	len  int
	sign int
}

func parseVarintPrefix(pre byte) *varintPrefixData {
	viPrefix := int(pre)
	viLen := 0
	sign := 0

	if viPrefix > 127 {
		viLen = viPrefix - 127
		sign = -1
	} else {
		viLen = viPrefix + 1
		sign = 1
	}

	return &varintPrefixData{
		len:  viLen,
		sign: sign,
	}
}

func parseVarint(buf []byte, startIdx int, curComponent *Component) (value int, readBytes int, err error) {
	// Interpret prefix

	if startIdx > len(buf)-1 {
		return 0, 0, errors.New("Could not complete parse. Incomplete Varint '" + curComponent.Name + "'.")
	}

	viPrefix := buf[startIdx]

	preData := parseVarintPrefix(viPrefix)

	if len(buf) < startIdx+preData.len+1 {
		return 0, 0, errors.New("Could not complete parse. Incomplete Varint '" + curComponent.Name + "'.")
	}

	viRaw := buf[startIdx+1 : startIdx+preData.len+1]
	viVal := 0

	for viIdx, viCurVal := range viRaw {
		placeVal := int(math.Pow(256, float64(viIdx)))

		viVal += placeVal * int(viCurVal)
	}

	viVal *= preData.sign

	// Return read data

	return viVal, preData.len + 1, nil
}

func buildVarint(ov interface{}) ([]byte, error) {
	value, ok := ov.(int)

	if !ok {
		return nil, errors.New("expected int as value for ov")
	}

	v := []byte{}

	osign := 1

	if value < 0 {
		osign = -1
	} else {
		osign = 1
	}

	value *= osign

	maxPlace := 1

	for value > maxPlace {
		maxPlace *= 256
	}

	maxPlace /= 256

	for maxPlace >= 1 {
		placeValue := int(math.Floor(float64(value) / float64(maxPlace)))
		v = append([]byte{byte(placeValue)}, v...)

		value %= maxPlace
		maxPlace /= 256
	}

	// Prepend length prefix

	viLen := len(v)

	if osign == -1 {
		// Negative number

		v = append([]byte{byte(127 + viLen)}, v...)
	} else {
		// Positive number

		v = append([]byte{byte(viLen - 1)}, v...)
	}

	return v, nil
}

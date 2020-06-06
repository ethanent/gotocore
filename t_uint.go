package gotocore

import (
	"encoding/binary"
	"errors"
	"github.com/ethanent/gotocore/util"
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

func buildUInt(value uint, size int) ([]byte, error) {
	if size < 1 {
		panic("invalid uint size " + strconv.Itoa(size))
	}

	build := make([]byte, size/8, size/8)

	if value > uint(util.Ipow(2, size)-1) {
		return nil, errors.New("value " + strconv.Itoa(int(value)) + " out of range for uint" + strconv.Itoa(size))
	}

	switch size {
	case 8:
		build[0] = byte(value)
	case 16:
		binary.LittleEndian.PutUint16(build, uint16(value))
	case 32:
		binary.LittleEndian.PutUint32(build, uint32(value))
	case 64:
		binary.LittleEndian.PutUint64(build, uint64(value))
	default:
		panic("invalid uint length " + strconv.Itoa(size))
	}

	return build, nil
}

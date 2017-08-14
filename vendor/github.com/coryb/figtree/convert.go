package figtree

import (
	"fmt"
	"strconv"
)

// dst must be a pointer type
func convertString(src string, dst interface{}) (err error) {
	switch v := dst.(type) {
	case *bool:
		*v, err = strconv.ParseBool(src)
	case *string:
		*v = src
	case *int:
		var tmp int64
		// this is a cheat, we only know int is at least 32 bits
		// but we have to make a compromise here
		tmp, err = strconv.ParseInt(src, 10, 32)
		*v = int(tmp)
	case *int8:
		var tmp int64
		tmp, err = strconv.ParseInt(src, 10, 8)
		*v = int8(tmp)
	case *int16:
		var tmp int64
		tmp, err = strconv.ParseInt(src, 10, 16)
		*v = int16(tmp)
	case *int32:
		var tmp int64
		tmp, err = strconv.ParseInt(src, 10, 32)
		*v = int32(tmp)
	case *int64:
		var tmp int64
		tmp, err = strconv.ParseInt(src, 10, 64)
		*v = int64(tmp)
	case *uint:
		var tmp uint64
		// this is a cheat, we only know uint is at least 32 bits
		// but we have to make a compromise here
		tmp, err = strconv.ParseUint(src, 10, 32)
		*v = uint(tmp)
	case *uint8:
		var tmp uint64
		tmp, err = strconv.ParseUint(src, 10, 8)
		*v = uint8(tmp)
	case *uint16:
		var tmp uint64
		tmp, err = strconv.ParseUint(src, 10, 16)
		*v = uint16(tmp)
	case *uint32:
		var tmp uint64
		tmp, err = strconv.ParseUint(src, 10, 32)
		*v = uint32(tmp)
	case *uint64:
		var tmp uint64
		tmp, err = strconv.ParseUint(src, 10, 64)
		*v = uint64(tmp)
	// hmm, collides with uint8
	// case *byte:
	// 	tmp := []byte(src)
	// 	if len(tmp) == 1 {
	// 		*v = tmp[0]
	// 	} else {
	// 		err = fmt.Errorf("Cannot convert string %q to byte, length: %d", src, len(tmp))
	// 	}
	// hmm, collides with int32
	// case *rune:
	// 	tmp := []rune(src)
	// 	if len(tmp) == 1 {
	// 		*v = tmp[0]
	// 	} else {
	// 		err = fmt.Errorf("Cannot convert string %q to rune, lengt: %d", src, len(tmp))
	// 	}
	case *float32:
		var tmp float64
		tmp, err = strconv.ParseFloat(src, 32)
		*v = float32(tmp)
	case *float64:
		var tmp float64
		tmp, err = strconv.ParseFloat(src, 64)
		*v = float64(tmp)
	default:
		err = fmt.Errorf("Cannot convert string %q to type %T", src, dst)
	}
	if err != nil {
		return err
	}

	return nil
}

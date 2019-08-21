package nullgo

import (
	"math"
	"unsafe"
)

func parseIntToBin(i int) (b []bool) {
	b = make([]bool, 8)
	pos := 7
	for i != 0 {
		add := false
		tmp := i % 2
		i = i / 2
		if tmp == 1 {
			add = true
		}
		b[pos] = add
		pos--
	}
	return
}

func parseBinToInt(b []bool) (res int) {
	pos := 0
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] {
			res += int(math.Pow(float64(2), float64(pos)))
		}
		pos++
	}
	return
}

func QuickStringToBytes(s string) (b []byte) {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func QuickBytesToString(b []byte) (s string)  {
	return *(*string)(unsafe.Pointer(&b))
}

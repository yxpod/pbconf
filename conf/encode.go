package conf

import "math"

func packInt(x int) []byte {
	var b []byte
	for x >= 1<<7 {
		b = append(b, uint8(x&0x7f|0x80))
		x >>= 7
	}
	b = append(b, uint8(x))
	return b
}

func packFieldInt(x int, n int) []byte {
	return append(packInt(n<<3), packInt(x)...)
}

func packFloat(x float32) []byte {
	n := math.Float32bits(x)
	return []byte{
		uint8(n),
		uint8(n >> 8),
		uint8(n >> 16),
		uint8(n >> 24),
	}
}

func packFieldFloat(x float32, n int) []byte {
	return append(packInt((n<<3)|5), packFloat(x)...)
}

func packBytes(b []byte) []byte {
	bb := packInt(len(b))
	return append(bb, b...)
}

func packFieldBytes(b []byte, n int) []byte {
	return append(packInt((n<<3)|2), packBytes(b)...)
}

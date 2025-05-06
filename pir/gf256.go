package pir

// Tiny GF(256) (AES poly 0x11b).  Enough for +,·,dot.

var exp = [512]byte{}
var logt = [256]byte{}

func init() {
	x := byte(1)
	for i := 0; i < 255; i++ {
		exp[i] = x
		logt[x] = byte(i)
		x ^= x << 1
		if int(x)&0x100 != 0 {
			x ^= 0x1b
		}
	}
	for i := 255; i < 512; i++ {
		exp[i] = exp[i-255]
	}
}

func Add(a, b byte) byte { return a ^ b }

func Mul(a, b byte) byte {
	if a == 0 || b == 0 {
		return 0
	}
	return exp[int(logt[a])+int(logt[b])]
}

// Dot product of equal‑length byte slices.
func Dot(a, b []byte) byte {
	var s byte
	for i := range a {
		s = Add(s, Mul(a[i], b[i]))
	}
	return s
}

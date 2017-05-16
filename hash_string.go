package bloom

// MurmurHash3 implementation adapted from Sébastien Paolacci
// github.com/spaolacci/murmur3, released under BSD-3-Clause.

func (d *digestString) hash(data string) (h1 uint64, h2 uint64) {
	d.h1, d.h2 = 0, 0
	d.clen = len(data)
	d.tail = d.bmixString(data)
	return d.sum()
}

type digestString struct {
	clen int
	tail string
	h1   uint64
	h2   uint64
}

func Uint64String(b string) uint64 {
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func (d *digestString) bmixString(p string) (tail string) {
	h1, h2 := d.h1, d.h2
	nblocks := len(p) / 16
	for i := 0; i < nblocks; i++ {
		j := 16 * i
		k1 := Uint64String(p[j : j+8])
		k2 := Uint64String(p[j+8 : j+16])
		k1 *= c1
		k1 = (k1 << 31) | (k1 >> 33)
		k1 *= c2
		h1 ^= k1
		h1 = (h1 << 27) | (h1 >> 37)
		h1 += h2
		h1 = h1*5 + 0x52dce729
		k2 *= c2
		k2 = (k2 << 33) | (k2 >> 31)
		k2 *= c1
		h2 ^= k2
		h2 = (h2 << 31) | (h2 >> 33)
		h2 += h1
		h2 = h2*5 + 0x38495ab5
	}
	d.h1, d.h2 = h1, h2
	return p[nblocks*16:]
}

func (d *digestString) sum() (h1, h2 uint64) {
	h1, h2 = d.h1, d.h2
	var k1, k2 uint64
	switch len(d.tail) & 15 {
	case 15:
		k2 ^= uint64(d.tail[14]) << 48
		fallthrough
	case 14:
		k2 ^= uint64(d.tail[13]) << 40
		fallthrough
	case 13:
		k2 ^= uint64(d.tail[12]) << 32
		fallthrough
	case 12:
		k2 ^= uint64(d.tail[11]) << 24
		fallthrough
	case 11:
		k2 ^= uint64(d.tail[10]) << 16
		fallthrough
	case 10:
		k2 ^= uint64(d.tail[9]) << 8
		fallthrough
	case 9:
		k2 ^= uint64(d.tail[8]) << 0
		k2 *= c2
		k2 = (k2 << 33) | (k2 >> 31)
		k2 *= c1
		h2 ^= k2
		fallthrough
	case 8:
		k1 ^= uint64(d.tail[7]) << 56
		fallthrough
	case 7:
		k1 ^= uint64(d.tail[6]) << 48
		fallthrough
	case 6:
		k1 ^= uint64(d.tail[5]) << 40
		fallthrough
	case 5:
		k1 ^= uint64(d.tail[4]) << 32
		fallthrough
	case 4:
		k1 ^= uint64(d.tail[3]) << 24
		fallthrough
	case 3:
		k1 ^= uint64(d.tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint64(d.tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint64(d.tail[0]) << 0
		k1 *= c1
		k1 = (k1 << 31) | (k1 >> 33)
		k1 *= c2
		h1 ^= k1
	}
	h1 ^= uint64(d.clen)
	h2 ^= uint64(d.clen)
	h1 += h2
	h2 += h1
	h1 = fmix(h1)
	h2 = fmix(h2)
	h1 += h2
	h2 += h1
	return h1, h2
}

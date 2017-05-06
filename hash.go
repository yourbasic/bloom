package bloom

import (
	"unsafe"
)

// MurmurHash3 implementation adapted from Sébastien Paolacci
// github.com/spaolacci/murmur3, released under BSD-3-Clause.

func (d *digest) hash(data []byte) (h1 uint64, h2 uint64) {
	d.h1, d.h2 = 0, 0
	d.clen = len(data)
	d.tail = d.bmix(data)
	return d.sum()
}

const (
	c1 = 0x87c37b91114253d5
	c2 = 0x4cf5ad432745937f
)

type digest struct {
	clen int      // Digested input cumulative length.
	buf  [16]byte // Expected (but not required) to be 16 large.
	tail []byte   // 0 to 15 bytes view of buf.
	h1   uint64   // Running hash part 1.
	h2   uint64   // Running hash part 2.
}

func (d *digest) bmix(p []byte) (tail []byte) {
	h1, h2 := d.h1, d.h2
	nblocks := len(p) / 16
	for i := 0; i < nblocks; i++ {
		t := (*[2]uint64)(unsafe.Pointer(&p[i*16]))
		k1, k2 := t[0], t[1]
		// Without unsafe on a little endian machine:
		// j := 16 * i
		// k1 := binary.LittleEndian.Uint64(p[j : j+8])
		// k2 := binary.LittleEndian.Uint64(p[j+8 : j+16])
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

func (d *digest) sum() (h1, h2 uint64) {
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

func fmix(k uint64) uint64 {
	k ^= k >> 33
	k *= 0xff51afd7ed558ccd
	k ^= k >> 33
	k *= 0xc4ceb9fe1a85ec53
	k ^= k >> 33
	return k
}

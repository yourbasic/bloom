package bloom

// MurmurHash3 implementation adapted from Sébastien Paolacci
// github.com/spaolacci/murmur3, released under BSD-3-Clause.

const (
	c1 = 0x87c37b91114253d5
	c2 = 0x4cf5ad432745937f
)

func fmix(k uint64) uint64 {
	k ^= k >> 33
	k *= 0xff51afd7ed558ccd
	k ^= k >> 33
	k *= 0xc4ceb9fe1a85ec53
	k ^= k >> 33
	return k
}

func uint64byte(b []byte) uint64 {
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func uint64string(s string) uint64 {
	return uint64(s[0]) | uint64(s[1])<<8 | uint64(s[2])<<16 | uint64(s[3])<<24 |
		uint64(s[4])<<32 | uint64(s[5])<<40 | uint64(s[6])<<48 | uint64(s[7])<<56
}

func hash(b []byte) (h1, h2 uint64) {
	nblocks := len(b) / 16
	for i := 0; i < nblocks; i++ {
		j := 16 * i
		k1, k2 := uint64byte(b[j:j+8]), uint64byte(b[j+8:j+16])

		k1 *= c1
		k1 = (k1 << 31) | (k1 >> 33) // rotl64(k1, 31)
		k1 *= c2

		h1 ^= k1
		h1 = (h1 << 27) | (h1 >> 37) // rotl64(h1, 27)
		h1 += h2
		h1 = h1*5 + 0x52dce729

		k2 *= c2
		k2 = (k2 << 33) | (k2 >> 31) // rotl64(k2, 33)
		k2 *= c1

		h2 ^= k2
		h2 = (h2 << 31) | (h2 >> 33) // rotl64(h2, 31)
		h2 += h1
		h2 = h2*5 + 0x38495ab5
	}

	tail := b[nblocks*16:]
	var k1, k2 uint64
	switch len(tail) {
	case 15:
		k2 ^= uint64(tail[14]) << 48
		fallthrough
	case 14:
		k2 ^= uint64(tail[13]) << 40
		fallthrough
	case 13:
		k2 ^= uint64(tail[12]) << 32
		fallthrough
	case 12:
		k2 ^= uint64(tail[11]) << 24
		fallthrough
	case 11:
		k2 ^= uint64(tail[10]) << 16
		fallthrough
	case 10:
		k2 ^= uint64(tail[9]) << 8
		fallthrough
	case 9:
		k2 ^= uint64(tail[8]) << 0
		k2 *= c2
		k2 = (k2 << 33) | (k2 >> 31) // rotl64(k2, 33)
		k2 *= c1
		h2 ^= k2
		fallthrough
	case 8:
		k1 ^= uint64(tail[7]) << 56
		fallthrough
	case 7:
		k1 ^= uint64(tail[6]) << 48
		fallthrough
	case 6:
		k1 ^= uint64(tail[5]) << 40
		fallthrough
	case 5:
		k1 ^= uint64(tail[4]) << 32
		fallthrough
	case 4:
		k1 ^= uint64(tail[3]) << 24
		fallthrough
	case 3:
		k1 ^= uint64(tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint64(tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint64(tail[0]) << 0
		k1 *= c1
		k1 = (k1 << 31) | (k1 >> 33) // rotl64(k1, 31)
		k1 *= c2
		h1 ^= k1
	}
	h1 ^= uint64(len(b))
	h2 ^= uint64(len(b))
	h1 += h2
	h2 += h1
	h1, h2 = fmix(h1), fmix(h2)
	h1 += h2
	h2 += h1
	return
}

func hashString(s string) (h1, h2 uint64) {
	nblocks := len(s) / 16
	for i := 0; i < nblocks; i++ {
		j := 16 * i
		k1, k2 := uint64string(s[j:j+8]), uint64string(s[j+8:j+16])

		k1 *= c1
		k1 = (k1 << 31) | (k1 >> 33) // rotl64(k1, 31)
		k1 *= c2

		h1 ^= k1
		h1 = (h1 << 27) | (h1 >> 37) // rotl64(h1, 27)
		h1 += h2
		h1 = h1*5 + 0x52dce729

		k2 *= c2
		k2 = (k2 << 33) | (k2 >> 31) // rotl64(k2, 33)
		k2 *= c1

		h2 ^= k2
		h2 = (h2 << 31) | (h2 >> 33) // rotl64(h2, 31)
		h2 += h1
		h2 = h2*5 + 0x38495ab5
	}

	tail := s[nblocks*16:]
	var k1, k2 uint64
	switch len(tail) {
	case 15:
		k2 ^= uint64(tail[14]) << 48
		fallthrough
	case 14:
		k2 ^= uint64(tail[13]) << 40
		fallthrough
	case 13:
		k2 ^= uint64(tail[12]) << 32
		fallthrough
	case 12:
		k2 ^= uint64(tail[11]) << 24
		fallthrough
	case 11:
		k2 ^= uint64(tail[10]) << 16
		fallthrough
	case 10:
		k2 ^= uint64(tail[9]) << 8
		fallthrough
	case 9:
		k2 ^= uint64(tail[8]) << 0
		k2 *= c2
		k2 = (k2 << 33) | (k2 >> 31) // rotl64(k2, 33)
		k2 *= c1
		h2 ^= k2
		fallthrough
	case 8:
		k1 ^= uint64(tail[7]) << 56
		fallthrough
	case 7:
		k1 ^= uint64(tail[6]) << 48
		fallthrough
	case 6:
		k1 ^= uint64(tail[5]) << 40
		fallthrough
	case 5:
		k1 ^= uint64(tail[4]) << 32
		fallthrough
	case 4:
		k1 ^= uint64(tail[3]) << 24
		fallthrough
	case 3:
		k1 ^= uint64(tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint64(tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint64(tail[0]) << 0
		k1 *= c1
		k1 = (k1 << 31) | (k1 >> 33) // rotl64(k1, 31)
		k1 *= c2
		h1 ^= k1
	}
	h1 ^= uint64(len(s))
	h2 ^= uint64(len(s))
	h1 += h2
	h2 += h1
	h1, h2 = fmix(h1), fmix(h2)
	h1 += h2
	h2 += h1
	return
}

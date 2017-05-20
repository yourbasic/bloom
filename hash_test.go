package bloom

import (
	"testing"
)

func TestHash(t *testing.T) {
	var data = []struct {
		h1, h2 uint64
		s      string
	}{
		{0x0000000000000000, 0x0000000000000000, ""},
		{0xcbd8a7b341bd9b02, 0x5b1e906a48ae1d19, "hello"},
		{0x342fac623a5ebc8e, 0x4cdcbc079642414d, "hello, world"},
		{0xb89e5988b737affc, 0x664fc2950231b2cb, "19 Jan 2038 at 3:14:07 AM"},
		{0xcd99481f9ee902c9, 0x695da1a38987b6e7, "The quick brown fox jumps over the lazy dog."},
	}
	for _, x := range data {
		h1, h2 := hash([]byte(x.s))
		if h1 != x.h1 {
			t.Errorf("hash(%q).h1 = %d; want %d\n", x.s, h1, x.h1)
		}
		if h2 != x.h2 {
			t.Errorf("hash(%q).h2 = %d; want %d\n", x.s, h2, x.h2)
		}
	}
}

func TestHashString(t *testing.T) {
	var data = []struct {
		h1, h2 uint64
		s      string
	}{
		{0x0000000000000000, 0x0000000000000000, ""},
		{0xcbd8a7b341bd9b02, 0x5b1e906a48ae1d19, "hello"},
		{0x342fac623a5ebc8e, 0x4cdcbc079642414d, "hello, world"},
		{0xb89e5988b737affc, 0x664fc2950231b2cb, "19 Jan 2038 at 3:14:07 AM"},
		{0xcd99481f9ee902c9, 0x695da1a38987b6e7, "The quick brown fox jumps over the lazy dog."},
	}
	for _, x := range data {
		h1, h2 := hashString(x.s)
		if h1 != x.h1 {
			t.Errorf("hashString(%q).h1 = %d; want %d\n", x.s, h1, x.h1)
		}
		if h2 != x.h2 {
			t.Errorf("hashString(%q).h2 = %d; want %d\n", x.s, h2, x.h2)
		}
	}
}

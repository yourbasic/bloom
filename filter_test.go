package bloom

import (
	"testing"
)

func TestFilter(t *testing.T) {
	s1 := "asöldkgjaösldkgaösldkasldgjkaösldkgjöasgkdjg"
	s2 := "elasödlnkgaölsdkfgaölsdkjfaölsdkgaölskgnaösl"
	s3 := "aölsdgkaösldkgaösldkgjaölsdkjgaölsdkgjaösldk"
	for n := 0; n < 100; n++ {
		for p := 1; p <= 128; p *= 2 {
			filter := New(n, p)
			member := filter.Likely(s1)
			if member {
				t.Errorf("Likely(s1) = %v; want false\n", member)
			}
			count := filter.Count()
			if count != 0 {
				t.Errorf("Count() = %d; want 0\n", count)
			}

			member = filter.Add(s1)
			if member {
				t.Errorf("Add(s1) = %v; want false\n", member)
			}
			count = filter.Count()
			if count != 1 {
				t.Errorf("Count() = %d; want 1\n", count)
			}
			member = filter.Likely(s1)
			if !member {
				t.Errorf("Likely(s1) = %v; want true\n", member)
			}
			member = filter.Likely(s2)
			if member {
				t.Errorf("Likely(s2) = %v; want false\n", member)
			}

			member = filter.Add(s1)
			if !member {
				t.Errorf("Add(s1) = %v; want true\n", member)
			}
			count = filter.Count()
			if count != 1 {
				t.Errorf("Count() = %d; want 1\n", count)
			}

			member = filter.Add(s3)
			if member {
				t.Errorf("Add(s1) = %v; want false\n", member)
			}
			count = filter.Count()
			if count != 2 {
				t.Errorf("Count() = %d; want 2\n", count)
			}
		}
	}
}

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
			member := filter.Test(s1)
			if member {
				t.Errorf("Test(s1) = %v; want false\n", member)
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
			member = filter.Test(s1)
			if !member {
				t.Errorf("Test(s1) = %v; want true\n", member)
			}
			member = filter.Test(s2)
			if member {
				t.Errorf("Test(s2) = %v; want false\n", member)
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
				t.Errorf("Add(s3) = %v; want false\n", member)
			}
			count = filter.Count()
			if count != 2 {
				t.Errorf("Count() = %d; want 2\n", count)
			}
		}
	}
}

func TestFilterByte(t *testing.T) {
	s1 := []byte("asöldkgjaösldkgaösldkasldgjkaösldkgjöasgkdjg")
	s2 := []byte("elasödlnkgaölsdkfgaölsdkjfaölsdkgaölskgnaösl")
	s3 := []byte("aölsdgkaösldkgaösldkgjaölsdkjgaölsdkgjaösldk")
	for n := 0; n < 100; n++ {
		for p := 1; p <= 128; p *= 2 {
			filter := New(n, p)
			member := filter.TestByte(s1)
			if member {
				t.Errorf("TestByte(s1) = %v; want false\n", member)
			}
			count := filter.Count()
			if count != 0 {
				t.Errorf("Count() = %d; want 0\n", count)
			}

			member = filter.AddByte(s1)
			if member {
				t.Errorf("AddByte(s1) = %v; want false\n", member)
			}
			count = filter.Count()
			if count != 1 {
				t.Errorf("Count() = %d; want 1\n", count)
			}
			member = filter.TestByte(s1)
			if !member {
				t.Errorf("TestByte(s1) = %v; want true\n", member)
			}
			member = filter.TestByte(s2)
			if member {
				t.Errorf("TestByte(s2) = %v; want false\n", member)
			}

			member = filter.AddByte(s1)
			if !member {
				t.Errorf("AddByte(s1) = %v; want true\n", member)
			}
			count = filter.Count()
			if count != 1 {
				t.Errorf("Count() = %d; want 1\n", count)
			}

			member = filter.AddByte(s3)
			if member {
				t.Errorf("AddByte(s3) = %v; want false\n", member)
			}
			count = filter.Count()
			if count != 2 {
				t.Errorf("Count() = %d; want 2\n", count)
			}
		}
	}
}

func TestUnion(t *testing.T) {
	s1 := "asöldkgjaösldkgaösldkasldgjkaösldkgjöasgkdjg"
	s2 := "elasödlnkgaölsdkfgaölsdkjfaölsdkgaölskgnaösl"
	s3 := "aölsdgkaösldkgaösldkgjaölsdkjgaölsdkgjaösldk"
	for n := 0; n < 100; n++ {
		for p := 1; p <= 128; p *= 2 {
			f1, f2 := New(n, p), New(n, p)
			f1.Add(s1)
			f1.Add(s2)
			f2.Add(s2)
			f2.Add(s3)
			or := f1.Union(f2)
			member := or.Test(s1)
			if !member {
				t.Errorf("f1.Union(f2).Test(s1) = %v; want true\n", member)
			}
			member = or.Test(s2)
			if !member {
				t.Errorf("f1.Union(f2).Test(s2) = %v; want true\n", member)
			}
			member = or.Test(s3)
			if !member {
				t.Errorf("f1.Union(f2).Test(s3) = %v; want true\n", member)
			}
		}
	}
}

var fox string = "The quick brown fox jumps over the lazy dog."

func BenchmarkAdd(b *testing.B) {
	b.StopTimer()
	filter := New(1<<30, 200)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = filter.Add(fox)
	}
}

func BenchmarkAddByte(b *testing.B) {
	b.StopTimer()
	filter := New(1<<30, 200)
	b.StartTimer()
	bytes := []byte(fox)
	for i := 0; i < b.N; i++ {
		_ = filter.AddByte(bytes)
	}
}

func BenchmarkTest(b *testing.B) {
	b.StopTimer()
	filter := New(1<<30, 200)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = filter.Test(fox)
	}
}

func BenchmarkTestByte(b *testing.B) {
	b.StopTimer()
	filter := New(1<<30, 200)
	b.StartTimer()
	bytes := []byte(fox)
	for i := 0; i < b.N; i++ {
		_ = filter.TestByte(bytes)
	}
}

func BenchmarkUnion(b *testing.B) {
	n := 1000
	b.StopTimer()
	f1 := New(n, 200)
	f2 := New(n, 200)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = f1.Union(f2)
	}
}

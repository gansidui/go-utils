package ringbuffer

import (
	"testing"
)

func TestRingBuffer(t *testing.T) {
	r := New(6)

	if r.Capacity() != 6 || r.Size() != 0 || r.FreeSize() != 6 {
		t.Fatal()
	}

	n := r.Write([]byte("love"))
	if n != 4 || r.Size() != 4 || r.FreeSize() != 2 || r.Capacity() != 6 {
		t.Fatal()
	}

	p := make([]byte, 10)
	n = r.Read(p, 10)
	if n != 0 {
		t.Fatal()
	}

	n = r.Read(p[2:], 3)
	if n != 3 || string(p[2:5]) != "lov" || r.Size() != 1 {
		t.Fatal()
	}

	n = r.Write([]byte("vent"))
	if n != 4 || r.Size() != 5 {
		t.Fatal()
	}

	n = r.Read(p, 5)
	if n != 5 || string(p[:5]) != "event" || r.Size() != 0 {
		t.Fatal()
	}
}

func TestRingBuffer2(t *testing.T) {
	r := New(11)

	r.Write([]byte("hello"))
	r.Reset()
	if r.Size() != 0 {
		t.Fatal()
	}

	r.Write([]byte("hello,world"))
	if r.Size() != 11 {
		t.Fatal()
	}

	r.RemoveNewest(6)

	p := make([]byte, 5)
	n := r.Read(p, 5)
	if n != 5 || string(p) != "hello" || r.Size() != 0 {
		t.Fatal()
	}

	r.Write([]byte("hello,world"))
	r.RemoveOldest(6)
	r.Read(p, 5)
	if string(p) != "world" || r.Size() != 0 {
		t.Fatal()
	}
}

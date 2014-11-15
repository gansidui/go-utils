package safemap

import (
	"testing"
)

func TestSafeMap(t *testing.T) {
	mp := New()

	if mp.Size() != 0 {
		t.Fatal()
	}

	mp.Set("hello", "world")

	if !mp.IsExist("hello") {
		t.Fatal()
	}

	mp.Set("golang", 100)
	mp.Set("google", 3.14)

	if val, ok := mp.Get("hello"); !ok || val.(string) != "world" {
		t.Fatal()
	}

	if val, ok := mp.Get("golang"); !ok || val.(int) != 100 {
		t.Fatal()
	}

	if val, ok := mp.Get("lijie"); ok || val != nil {
		t.Fatal()
	}

	if mp.Size() != 3 {
		t.Fatal()
	}

	mp.Delete("hello")

	m := mp.Items()
	if len(m) != 2 || m["golang"].(int) != 100 || m["google"].(float64) != 3.14 {
		t.Fatal()
	}

}

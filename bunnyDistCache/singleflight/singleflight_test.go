package singleflight

import "testing"

func TestDo(t *testing.T) {
	var f Flight
	v, err := f.Fly("key", func() (interface{}, error) {
		return "bar", nil
	})

	if v != "bar" || err != nil {
		t.Errorf("Do v = %v, error = %v", v, err)
	}
}

package input

import "testing"

func TestDef(t *testing.T) {
	i := NewInputter()
	if err := i.Test(); err != nil {
		t.Fatal(err)
	}
}

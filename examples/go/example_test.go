package main

import "testing"

func TestSum(t *testing.T) {
	if sum(1, 2) != 3 {
		t.Errorf("sum(1, 2) != 3")
	}

	if sum(10, 20) != 30 {
		t.Errorf("sum(10, 20) != 30")
	}
}

package main

import "testing"

func TestHighSum(t *testing.T) {
	if sum(1000, 1000) != 2000 {
		t.Errorf("sum(1000, 1000) != 2000")
	}
}

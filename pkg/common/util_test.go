package common

import "testing"

func TestUnitFloat2Str(t *testing.T) {
	if Float2Str(123.9921) != "123.9921" {
		t.Fatal(Float2Str(123.9921), "123.9921")
	}
	if Float2Str(0.00000111205) != "0.00000111205" {
		t.Fatal(Float2Str(0.00000111205), "0.00000111205")
	}
	if Float2Str(112312.00) != "112312" {
		t.Fatal(Float2Str(112312.00), "112312")
	}
	if Float2Str(0.0000000000000228) != "0.0000000000000228" {
		t.Fatal(Float2Str(0.0000000000000228), "0.0000000000000228")
	}
	if Float2Str(112312.01129210) != "112312.0112921" {
		t.Fatal(Float2Str(112312.01129210), "112312.0112921")
	}
}

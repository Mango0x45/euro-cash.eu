/* This file contains tests for utility functions used in the mfmt
   binary.  Tests of the actual application are in main_test.go. */

package main

import (
	"testing"

	"git.thomasvoss.com/euro-cash.eu/mintage"
)

func TestFormatInt(t *testing.T) {
	for _, x := range [...]struct {
		n int
		s string
	}{
		{0, "0"},
		{123, "123"},
		{81758, "81.758"},
		{752759237528, "752.759.237.528"},
		{mintage.Unknown, "?"},
	} {
		s := formatInt(x.n)
		if s != x.s {
			t.Fatalf(`Expected s="%s"; got "%s"`, x.s, s)
		}
	}
}

func TestIntLen(t *testing.T) {
	for _, x := range [...]struct{ x, y int }{
		{0, len("0")},
		{123, len("123")},
		{81758, len("81758")},
		{752759237528, len("752759237528")},
		{mintage.Unknown, len("?")},
	} {
		n := intlen(x.x)
		if n != x.y {
			t.Fatalf("Expected n=%d; got %d", x.y, n)
		}
	}
}

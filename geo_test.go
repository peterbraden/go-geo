package geo

import "testing"

func TestSmoke(t *testing.T) {
	var b1 = BBox{W: 25., E: 28., N: 10., S: 5.}
	var b2 = BBox{W: 26., E: 27., N: 8., S: 6.}

	if !b1.ContainsBBox(b2) {
		t.Fail()
	}

	if !b1.Contains(6., 26.) {
		t.Fail()
	}
}

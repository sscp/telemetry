package sundae

import (
	"context"
	"math"
	"testing"
)

type verifyFloat32Test struct {
	inVal  *float32
	outVal *float32
}

func (test verifyFloat32Test) runTest(t *testing.T) {
	ctx := context.Background()
	output := verifyFloat32(ctx, test.inVal, "testValue")
	if test.outVal != output {
		t.Errorf("error: verifyFloat32(%v) != %v (should be %v)",
			test.inVal, output, test.outVal)
	}
}

var positiveinf = float32(math.Inf(1))
var negativeinf = float32(math.Inf(-1))
var nan = float32(math.NaN())
var zero = float32(0.0)
var posnum = float32(10)
var negnum = float32(-10)

var float32Tests = []verifyFloat32Test{
	verifyFloat32Test{inVal: &posnum, outVal: &posnum},
	verifyFloat32Test{inVal: &negnum, outVal: &negnum},
	verifyFloat32Test{inVal: &zero, outVal: &zero},

	verifyFloat32Test{inVal: &positiveinf, outVal: nil},
	verifyFloat32Test{inVal: &negativeinf, outVal: nil},
	verifyFloat32Test{inVal: &nan, outVal: nil},
}

func TestVerifyFloat32(t *testing.T) {
	for _, test := range float32Tests {
		test.runTest(t)
	}
}

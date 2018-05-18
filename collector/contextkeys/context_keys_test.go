package contextkeys

import (
	"context"
	"testing"
	"time"
)

func TestContextValues(t *testing.T) {
	ctx := context.TODO()
	tm, ok := RecievedTimeFromContext(ctx)
	if ok {
		t.Errorf("Context (%v) without time should not be ok", ctx)
	}
	zerotm := time.Time{}
	if !ok && tm != zerotm {
		t.Errorf("Context (%v) without time should return zero time, but instead returned this: %v", ctx, tm)
	}
	nowTime := time.Now()
	ctx = ContextWithRecievedTime(ctx, nowTime)
	tm, ok = RecievedTimeFromContext(ctx)
	if ok && tm != nowTime {
		t.Errorf("Context (%v) time (%v) should equal this: %v", ctx, tm, nowTime)
	}
	if !ok && tm != zerotm {
		t.Errorf("Context (%v) without time should return zero time, but instead returned this: %v", ctx, tm)
	}

}

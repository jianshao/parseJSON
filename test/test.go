package parseJSON

import "testing"

func AssertEqual(t *testing.T, a interface{}, b interface{}, extra string)  {
	t.Helper()
	if a != b {
		t.Errorf("not equal, expect %v but %v, get %s", b, a, extra)
	}
}

func AssertNotEqual(t *testing.T, a interface{}, b interface{})  {
	t.Helper()
	if a == b {
		t.Errorf("equal")
	}
}

func AssertTrue(t *testing.T, a interface{})  {
	t.Helper()
	if a != true {
		t.Errorf("not true")
	}
}

func AssertFalse(t *testing.T, a interface{})  {
	t.Helper()
	if a != false {
		t.Errorf("not false")
	}
}

func AssertNil(t *testing.T, a interface{})  {
	t.Helper()
	if a != nil {
		t.Errorf("not nil, expect nil but %v", a)
	}
}

func AssertNotNil(t *testing.T, a interface{})  {
	t.Helper()
	if a == nil {
		t.Errorf("nil")
	}
}
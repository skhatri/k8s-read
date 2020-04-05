package test

import "testing"

func NotNull(t *testing.T, in interface{}) {
	if in == nil {
		t.Fail()
	}
}

func Null(t *testing.T, in interface{}) {
	if in != nil {
		t.Fail()
	}
}

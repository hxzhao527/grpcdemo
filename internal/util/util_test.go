package util

import (
	"testing"
)

func TestGetServiceNameFromFullMethod(t *testing.T) {
	t.Log(GetServiceNameFromFullMethod("/a.v/b"))
}

func TestGetIPFromAddress(t *testing.T) {
	t.Log(GetSelfIPAddress().String())
}

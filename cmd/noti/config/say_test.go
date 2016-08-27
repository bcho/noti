// +build darwin

package config

import (
	"reflect"
	"testing"

	"github.com/variadico/noti/speech"
)

func TestMergeSay(t *testing.T) {
	n1 := &say.Notification{
		Message: "hello",
		Voice:   "world",
		Rate:    10,
	}

	n2 := &say.Notification{
		Message: "good",
		Voice:   "bye",
		Rate:    20,
	}

	n3 := &say.Notification{
		Message: "fizz",
		Voice:   "buzz",
		Rate:    202,
	}

	want := &say.Notification{
		Message: "fizz",
		Voice:   "buzz",
		Rate:    202,
	}

	got := MergeSay(n1, n2, n3)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: '%v'; want: '%v'", got, want)
	}
}

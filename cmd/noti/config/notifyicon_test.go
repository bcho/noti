// +build !darwin
// +build !windows

package config

import (
	"reflect"
	"testing"

	"github.com/variadico/noti/freedesktop"
)

func TestMergeFreedesktop(t *testing.T) {
	n1 := &freedesktop.Notification{
		AppName:    "hello",
		ReplacesID: 10,
		AppIcon:    "world",
	}

	n2 := &freedesktop.Notification{
		Summary: "foo",
		Body:    "bar",
		Actions: []string{"fizz", "buzz"},
	}

	n3 := &freedesktop.Notification{
		Hints:         map[string]freedesktop.Hint{},
		ExpireTimeout: 100,
	}

	want := &freedesktop.Notification{
		AppName:       "hello",
		ReplacesID:    10,
		AppIcon:       "world",
		Summary:       "foo",
		Body:          "bar",
		Actions:       []string{"fizz", "buzz"},
		Hints:         map[string]freedesktop.Hint{},
		ExpireTimeout: 100,
	}

	got := MergeFreedesktop(n1, n2, n3)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: '%v'; want: '%v'", got, want)
	}
}

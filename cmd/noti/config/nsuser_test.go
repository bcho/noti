// +build darwin

package config

import (
	"reflect"
	"testing"

	"github.com/variadico/noti/nsuser"
)

func TestMergeNSUser(t *testing.T) {
	n1 := &nsuser.Notification{
		Title:           "hello",
		InformativeText: "world",
	}

	n2 := &nsuser.Notification{
		SoundName:    "foo",
		ContentImage: "bar.jpg",
	}

	n3 := &nsuser.Notification{
		Title: "testing",
	}

	want := &nsuser.Notification{
		Title:           "testing",
		Subtitle:        "",
		InformativeText: "world",
		ContentImage:    "bar.jpg",
		SoundName:       "foo",
	}

	got := MergeNSUser(n1, n2, n3)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: '%v'; want: '%v'", got, want)
	}
}

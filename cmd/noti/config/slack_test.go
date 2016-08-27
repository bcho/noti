package config

import (
	"reflect"
	"testing"

	"github.com/variadico/noti/slack"
)

func TestMergeSlack(t *testing.T) {
	n1 := &slack.Notification{
		Token:   "foo",
		Channel: "bar",
		Text:    "fizz",
		Parse:   "buzz",
	}

	n2 := &slack.Notification{
		LinkNames:   1,
		Attachments: nil,
		UnfurlLinks: true,
		UnfurlMedia: true,
	}

	n3 := &slack.Notification{
		Username:  "test",
		AsUser:    true,
		IconURL:   "abc",
		IconEmoji: "def",
	}

	want := &slack.Notification{
		Token:       "foo",
		Channel:     "bar",
		Text:        "fizz",
		Parse:       "buzz",
		LinkNames:   1,
		Attachments: nil,
		UnfurlLinks: true,
		UnfurlMedia: true,
		Username:    "test",
		AsUser:      true,
		IconURL:     "abc",
		IconEmoji:   "def",
	}

	got := MergeSlack(n1, n2, n3)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: '%v'; want: '%v'", got, want)
	}
}
